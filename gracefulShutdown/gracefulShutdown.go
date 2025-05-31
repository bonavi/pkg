package gracefulShutdown

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const gracefulShutdownDelay = 2 * time.Second

func AddGracefulShutdownErrGroup(
	serversErrWg *errgroup.Group,
	ctx context.Context,
	httpServers []*fiber.App,
	grpcServers []*grpc.Server,
) {

	// Создаем горутину, которая слушает основной контекст приложения, который завершается по сигналу от ОС +
	// Контекст от errgroup, который отменяется, когда хоть одна из горутин возвращает ошибку (или все горутины завершены, что невозможно в этом случае)
	serversErrWg.Go(func() error {

		<-ctx.Done()

		// Задерживаемся на пару секунд, чтобы если завершились с ошибкой/по kill при запуске сервиса, то успели стартануть сервера,
		// Иначе если до запуска будет вызван shutdown, то fiber не выдаст ошибку на Serve() и не сможем корректно завершить работу
		time.Sleep(gracefulShutdownDelay)

		// Плавно завершаем работу серверов
		for _, httpServer := range httpServers {
			_ = httpServer.ShutdownWithContext(ctx)
		}
		for _, grpcServer := range grpcServers {
			grpcServer.GracefulStop()
		}

		// Если мы до сюда дошли, значит либо одна из горутин вернула ошибку, либо контекст завершился по сигналу от ОС
		// В первом случае errgroup и так вернет ошибку на Wait(), во втором случае это обычное поведение на завершение работы
		return nil
	})

}
