package universesdedata

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go-updater/sde/loadyaml"
)

const (
	AFFIX_UNIVERSE = "fsd/universe"
)

type UniverseSDEData struct {
	ETCORegions map[b.RegionId]b.RegionName
	ETCOSystems map[b.SystemId]b.System
}

func TransceiveLoadAndConvert(
	pathSDE string,
	chnSend chanresult.ChanSendResult[UniverseSDEData],
) error {
	etcoUniverseSDEData, err := LoadAndConvert(pathSDE)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(etcoUniverseSDEData)
	}
}

func LoadAndConvert(
	pathSDE string,
) (
	etcoUniverseSDEData UniverseSDEData,
	err error,
) {
	regions := make(map[b.RegionId]b.RegionName)
	systems := make(map[b.SystemId]b.System)
	if err := getRegionsAndSystems(
		fmt.Sprintf("%s/%s", pathSDE, AFFIX_UNIVERSE),
		regions,
		systems,
	); err != nil {
		return etcoUniverseSDEData, err
	} else {
		return UniverseSDEData{
			ETCORegions: regions,
			ETCOSystems: systems,
		}, nil
	}
}

func getRegionsAndSystems(
	universeRoot string,
	regions map[b.RegionId]b.RegionName,
	systems map[b.SystemId]b.System,
) error {
	return filepath.Walk(universeRoot, func(
		path string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return err
		} else if isRegion(info) {
			id, val, err := toRegion(path, info)
			if err != nil {
				return err
			}
			regions[id] = val
			return getSystems(filepath.Dir(path), id, systems)
		}
		return nil
	})
}

func getSystems(
	regionRoot string,
	regionId int32,
	systems map[b.SystemId]b.System,
) error {
	return filepath.Walk(regionRoot, func(
		path string,
		info os.FileInfo,
		err error,
	) error {
		if err != nil {
			return err
		}

		if isSystem(info) {
			id, val, err := toSystem(path, info, regionId)
			if err != nil {
				return err
			}
			systems[id] = val
		}
		return nil
	})
}

func isRegion(info os.FileInfo) bool {
	return !info.IsDir() && info.Name() == "region.staticdata"
}

func isSystem(info os.FileInfo) bool {
	return !info.IsDir() && info.Name() == "solarsystem.staticdata"
}

func toRegion(
	path string,
	info os.FileInfo,
) (id int32, val b.RegionName, err error) {
	ymlData, err := loadyaml.LoadYaml[RegionStaticData](path)
	if err != nil {
		return id, val, err
	}
	regionName := fixRegionName(filepath.Base(filepath.Dir(path)))
	return ymlData.RegionId, regionName, nil
}

func toSystem(
	path string,
	info os.FileInfo,
	regionId int32,
) (id int32, val b.System, err error) {
	ymlData, err := loadyaml.LoadYaml[SolarSystemStaticData](path)
	if err != nil {
		return id, val, err
	}
	systemName := filepath.Base(filepath.Dir(path))
	return ymlData.SolarSystemId, b.System{
		RegionId: regionId,
		Name:     systemName,
	}, nil
}
