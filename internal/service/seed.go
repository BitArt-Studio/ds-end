package service

import (
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/model"
	"gohub/pkg/hashidsP"
	"gohub/pkg/snowflakeP"
	"gorm.io/gorm"
)

type SeedService struct {
}

var seedDao = dao.Seed
var Seed = new(SeedService)

// address => hSeed
var usedTempSeeds = make(map[string]string)

func (s *SeedService) RandomUsableSeed(address string) (string, error) {
	hSeed, err := hashidsP.HashID.EncodeInt64([]int64{snowflakeP.Node.Generate().Int64()})
	if err != nil {
		return "", errors.WithStack(err)
	}
	usedTempSeeds[address] = hSeed

	return hSeed, nil
}

func (s *SeedService) useSeed(tx *gorm.DB, hSeed string, address string) error {
	address = dealAddress(address)

	if err := seedDao.Tx(tx).New().Create(&model.SeedDO{
		Address: address,
		HSeed:   hSeed,
	}).Error; err != nil {
		return errors.WithStack(err)
	}

	delete(usedTempSeeds, address)
	return nil
}

func (s *SeedService) UsedTempSeed(address string) string {
	return usedTempSeeds[address]
}

func (s *SeedService) GetSeedsByAddress(address string) ([]string, error) {
	address = dealAddress(address)
	hSeeds := make([]string, 0)
	if err := seedDao.Model().Select("hSeed").Where("address = ?", address).Find(&hSeeds).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return hSeeds, nil
}
