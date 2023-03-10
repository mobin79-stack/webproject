package boot

import (
	"log"

	"git.ecobin.ir/ecomicro/bootstrap/service"
	book "git.ecobin.ir/ecomicro/template/app/book/boot"
	category "git.ecobin.ir/ecomicro/template/app/category/boot"
	comment "git.ecobin.ir/ecomicro/template/app/comment/boot"
	favorite "git.ecobin.ir/ecomicro/template/app/favorite/boot"
	purchase "git.ecobin.ir/ecomicro/template/app/purchase/boot"
	user "git.ecobin.ir/ecomicro/template/app/user/boot"
	"git.ecobin.ir/ecomicro/transport"
	"git.ecobin.ir/ecomicro/x/structure"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Boot(service *service.Service) {

	t, err := transport.NewTransport(service)
	if err != nil {
		log.Fatal(err, "Failed to create new transport!")
	}

	boots := make([]structure.BootInterface, 0)
	// boot domains
	boots = append(boots, user.Boot(service, t))
	boots = append(boots, book.Boot(service, t))
	boots = append(boots, comment.Boot(service, t))
	boots = append(boots, purchase.Boot(service, t))
	boots = append(boots, category.Boot(service, t))
	boots = append(boots, favorite.Boot(service, t))
	//TODO: add other domains
	bootData := structure.Boot{
		GrpcServers:  make(map[string]interface{}),
		GrpcClients:  make(map[string]interface{}),
		Usecases:     make(map[string]interface{}),
		Repositories: make(map[string]interface{}),
		Adapters:     make(map[string]interface{}),
	}
	_, err = t.HTTP("main", func(g *gin.Engine) { bootData.Gin = g })

	// repository
	for _, boot := range boots {
		boot.ApplyRepository(bootData)
	}
	// usecase
	for _, boot := range boots {
		boot.ApplyUsecase(bootData)
	}
	// http
	for _, boot := range boots {
		boot.ApplyHttpHandler(bootData)
	}
	// adapter
	for _, boot := range boots {
		boot.ApplyAdapters(bootData)
	}
	// swagger
	bootData.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
