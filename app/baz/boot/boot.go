package boot

import (
	"fmt"
	"log"

	"git.ecobin.ir/ecomicro/bootstrap/service"
	"git.ecobin.ir/ecomicro/template/app/baz/domain"
	bazRepo "git.ecobin.ir/ecomicro/template/app/baz/repository"
	bazUsecase "git.ecobin.ir/ecomicro/template/app/baz/usecase"
	mainDomain "git.ecobin.ir/ecomicro/template/domain"
	"git.ecobin.ir/ecomicro/transport"
	"gorm.io/gorm"
)

func getFromMap[T any](stringMap map[string]interface{}, key string) T {
	if value, ok := stringMap[key]; ok {
		if v, ok := value.(T); ok {
			return v
		}
		panic(fmt.Sprintf("assertion failed: %+v", stringMap[key]))
	}
	panic(fmt.Sprintf("key not found: map is=> %+v  and key is %+v", stringMap, key))
}

type bazBoot struct {
	transport *transport.Transport
	service   *service.Service
	db        *gorm.DB
}

var _ mainDomain.Boot = &bazBoot{}

func (b *bazBoot) ApplyRepository(boot mainDomain.DomainBoot) {
	if _, ok := boot.Repositories["baz"]; ok {
		log.Fatalf("baz repository already exist in repository map.")
	}
	boot.Repositories["baz"] = bazRepo.NewBazRepository(b.db)
}
func (b *bazBoot) ApplyUsecase(boot mainDomain.DomainBoot) {
	bazRepository := getFromMap[domain.Repository](boot.Repositories, "baz")
	if _, ok := boot.Usecases["baz"]; ok {
		log.Fatalf("baz usecase already exist in usecase map.")
	}
	boot.Usecases["baz"] = bazUsecase.NewBazUsecase(bazRepository)
}
func (b *bazBoot) ApplyHttpHandler(boot mainDomain.DomainBoot) {}
func (b *bazBoot) ApplyGrpcHandler(boot mainDomain.DomainBoot) {}
func (b *bazBoot) ApplyAdapters(boot mainDomain.DomainBoot)    {}

func Boot(service *service.Service, transport *transport.Transport) *bazBoot {
	return &bazBoot{
		transport: transport,
		service:   service,
		db:        service.Database["template"].GormDB,
	}
}
