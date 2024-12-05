package util

import (
	"gf2gacha/config"
	"gf2gacha/logger"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

const (
	GameDirNameCN     = "少女前线2：追放"
	GameDirNameGlobal = "EXILIUM"
)

func GetGameDataDir() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}

	var logData []byte
	logPathCN := filepath.Join(userHome, "/AppData/LocalLow/SunBorn/", GameDirNameCN, "/Player.log")
	logPathGlobal := filepath.Join(userHome, "/AppData/LocalLow/SunBorn/", GameDirNameGlobal, "/Player.log")
	if data, err := os.ReadFile(logPathCN); err == nil {
		logData = data
	} else if data, err := os.ReadFile(logPathGlobal); err == nil {
		logData = data
	} else {
		return "", errors.WithStack(err)
	}

	regexpGamePath, err := regexp.Compile(`\[Subsystems] Discovering subsystems at path (.+)/UnitySubsystems`)
	if err != nil {
		return "", errors.WithStack(err)
	}
	resultGamePath := regexpGamePath.FindSubmatch(logData)
	if len(resultGamePath) == 2 {
		gameDataDir := filepath.Join(string(resultGamePath[1]), "LocalCache/Data")

		//保存游戏路径，供日志里找不到时降级使用
		err = config.SetGameDataDir(gameDataDir)
		if err != nil {
			return "", err
		}

		return gameDataDir, nil
	}

	logger.Logger.Warn("未在日志中找到游戏路径,尝试从配置文件读取")
	//如果日志里没找到游戏路径，尝试降级读取配置文件
	if !config.IsSetGameDataDir() {
		return "", errors.New("日志和配置文件中均未找到游戏路径")
	}

	return config.GetGameDataDir(), nil
}
