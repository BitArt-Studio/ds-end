package service

import (
	"bytes"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
	"gohub/internal/dao"
	"gohub/internal/model"
	"gohub/pkg/hashidsP"
	"gohub/pkg/snowflakeP"
	"gorm.io/gorm"
	"html/template"
	"math"
	"sync"
)

type SeedService struct {
}

var seedDao = dao.Seed
var Seed = new(SeedService)
var mu sync.Mutex

// address => hSeed
var usedTempSeeds = make(map[string]string)

func (s *SeedService) generateShortID(node *snowflake.Node) int64 {
	snowflakeID := node.Generate().Int64()
	jsMaxInt := int64(math.Pow(2, 53) - 1)
	shortID := snowflakeID % jsMaxInt
	return shortID
}

func (s *SeedService) RandomUsableSeed(address string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	address = dealAddress(address)
	hSeed, err := hashidsP.HashID.EncodeInt64([]int64{s.generateShortID(snowflakeP.Node)})
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
	address = dealAddress(address)
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

func (s *SeedService) FillTemplate(filePath, hSeed string) ([]byte, error) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data := struct {
		HSeed string
	}{
		HSeed: hSeed,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}
