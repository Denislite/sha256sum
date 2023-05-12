package initialize

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/integrity-sum/internal/core/services"
	"github.com/integrity-sum/internal/repositories"
	"github.com/sirupsen/logrus"
)

func Initialize(ctx context.Context, logger *logrus.Logger, sig chan os.Signal) {
	// Initialize repository
	repository := repositories.NewAppRepository(logger)

	// Initialize service
	algorithm := os.Getenv("ALGORITHM")

	service := services.NewAppService(repository, algorithm, logger)

	// Initialize kubernetesAPI
	dataFromK8sAPI, err := service.GetDataFromK8sAPI()
	if err != nil {
		logger.Fatalf("can't get data from K8sAPI: %s", err)
	}

	//Getting pid
	pid, err := service.GetPID(dataFromK8sAPI.ConfigMapData)
	if err != nil {
		logger.Fatalf("err while getting pid %s", err)
	}
	if pid == 0 {
		logger.Fatalf("proc with name %s not exist", dataFromK8sAPI.ConfigMapData.ProcName)
	}

	//Getting the path to the monitoring directory
	dirPath := "../proc/" + strconv.Itoa(pid) + "/root/" + dataFromK8sAPI.ConfigMapData.MountPath

	duration, err := strconv.Atoi(os.Getenv("DURATION_TIME"))
	if err != nil {
		duration = 15
	}
	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context, ticker *time.Ticker) {
		defer wg.Done()
		for {
			if service.IsExistDeploymentNameInDB(dataFromK8sAPI.KuberData.TargetName) {
				logger.Info("Deployment name does not exist in database, save data")
				err := service.Start(ctx, dirPath, sig, dataFromK8sAPI.DeploymentData)
				if err != nil {
					logger.Fatalf("Error when starting to get and save hash data %s", err)
				}
			} else {
				logger.Info("Deployment name exists in database, checking data")
				for range ticker.C {
					err := service.Check(ctx, dirPath, sig, dataFromK8sAPI.DeploymentData, dataFromK8sAPI.KuberData)
					if err != nil {
						logger.Fatalf("Error when starting to check hash data %s", err)
					}
					logger.Info("Check completed")
				}
			}
		}
	}(ctx, ticker)
	wg.Wait()
	ticker.Stop()
}
