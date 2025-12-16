package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"

	"github.com/iviv660/wb-ImageProcessor/internal/adapter"
	consumerworker "github.com/iviv660/wb-ImageProcessor/internal/adapter/consumer"
	v1 "github.com/iviv660/wb-ImageProcessor/internal/api/v1"
	imagerepo "github.com/iviv660/wb-ImageProcessor/internal/repository/image"
	service "github.com/iviv660/wb-ImageProcessor/internal/service/image"
	localstorage "github.com/iviv660/wb-ImageProcessor/internal/storage/image"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpAddr := mustEnv("HTTP_ADDR")
	pgDSN := mustEnv("PG_DSN")

	kafkaBrokers := splitCSV(mustEnv("KAFKA_BROKERS"))
	kafkaTopic := mustEnv("KAFKA_TOPIC")
	kafkaGroupID := mustEnv("KAFKA_GROUP_ID")

	uploadsDir := mustEnv("UPLOADS_DIR")
	filesBaseURL := mustEnv("FILES_BASE_URL")

	db, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := imagerepo.New(db)
	st := localstorage.New(uploadsDir, filesBaseURL)

	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaBrokers...),
		Topic:        kafkaTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}
	defer writer.Close()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kafkaBrokers,
		Topic:          kafkaTopic,
		GroupID:        kafkaGroupID,
		CommitInterval: 0,
	})
	defer reader.Close()

	prod := NewKafkaProducer(writer)
	cons := NewKafkaConsumer(reader)

	svc := service.New(repo, st, prod)

	worker := consumerworker.NewWorker(cons, svc)
	go func() {
		if err := worker.Run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("worker stopped: %v", err)
			stop()
		}
	}()

	mux := chi.NewRouter()

	// Раздача файлов, если твой LocalStorage возвращает URL вида http://.../files/<name>
	mux.Handle("/files/*", http.StripPrefix("/files/", http.FileServer(http.Dir(uploadsDir))))

	api := v1.New(mux, svc, st)
	api.Register()

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("http listening on %s", httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func mustEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatalf("missing env %s", key)
	}
	return v
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

type KafkaProducer struct{ w *kafka.Writer }

func NewKafkaProducer(w *kafka.Writer) adapter.Producer { return &KafkaProducer{w: w} }

func (p *KafkaProducer) Send(ctx context.Context, key, value []byte) error {
	return p.w.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

type KafkaConsumer struct{ r *kafka.Reader }

func NewKafkaConsumer(r *kafka.Reader) adapter.Consumer { return &KafkaConsumer{r: r} }

func (c *KafkaConsumer) Consume(ctx context.Context, handle func(ctx context.Context, msg kafka.Message) error) error {
	for {
		m, err := c.r.FetchMessage(ctx)
		if err != nil {
			return err
		}
		if err := handle(ctx, m); err != nil {
			continue
		}
		if err := c.r.CommitMessages(ctx, m); err != nil {
			continue
		}
	}
}
