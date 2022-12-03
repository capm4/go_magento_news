package registry

import (
	"magento/bot/pkg/config"
	"magento/bot/pkg/controller"
	"magento/bot/pkg/database"
	"magento/bot/pkg/repository"
)

type Registry struct {
	Config         config.Сonfig
	dbConn         database.PostgresDB
	WebRepository  repository.WebsiteRepositoryInterface
	CfgRepository  repository.ConfigRepositoryInterface
	UserRepository repository.UserRepositoryInterface
}

func Init(cfg *config.Сonfig) (*Registry, error) {
	registry := Registry{}
	dbConn, err := CreateDBConnection(cfg)
	if err != nil {
		return nil, err
	}
	registry.dbConn = *dbConn
	webRep, err := CreateWebRepository(registry.dbConn)
	if err != nil {
		return nil, err
	}
	registry.WebRepository = webRep
	cfgRep, err := CreateCfgRepository(registry.dbConn)
	if err != nil {
		return nil, err
	}
	registry.CfgRepository = cfgRep
	usrRep, err := CreateUserRepository(registry.dbConn)
	if err != nil {
		return nil, err
	}
	registry.UserRepository = usrRep
	registry.Config = *cfg

	return &registry, nil
}

func CreateDBConnection(cfg *config.Сonfig) (*database.PostgresDB, error) {
	return database.NewPostgresDB(cfg.DBhost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBname)
}

//create website repository
func CreateWebRepository(dbCon database.PostgresDB) (repository.WebsiteRepositoryInterface, error) {
	db, err := database.NewPostgresWebsiteDB(&dbCon)
	if err != nil {
		return nil, err
	}
	return repository.NewWebsiteRepository(db), nil
}

//create config repository
func CreateCfgRepository(dbCon database.PostgresDB) (repository.ConfigRepositoryInterface, error) {
	db, err := database.NewPostgresConfigDB(&dbCon)
	if err != nil {
		return nil, err
	}
	return repository.NewConfigRepository(db), nil
}

//create user repository
func CreateUserRepository(dbCon database.PostgresDB) (repository.UserRepositoryInterface, error) {
	db, err := database.NewPostgresUserDB(&dbCon)
	if err != nil {
		return nil, err
	}
	return repository.NewUserRepository(db), nil
}

func (r Registry) NewAppController() (controller.AppController, error) {
	webController, err := controller.NewWebsiteController(r.WebRepository)
	if err != nil {
		return controller.AppController{}, err
	}
	userController, err := controller.NewUserController(r.UserRepository)
	if err != nil {
		return controller.AppController{}, err
	}
	return controller.AppController{Website: webController, User: userController, Config: &r.Config}, nil
}
