package edp_component

import (
	"edp-admin-console/context"
	"edp-admin-console/models/query"
	ec "edp-admin-console/repository/edp-component"
	"edp-admin-console/service/logger"
	dberror "edp-admin-console/util/error/db-errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

var log = logger.GetLogger()

type EDPComponentService struct {
	IEDPComponent ec.IEDPComponentRepository
}

//GetEDPComponent gets EDP component by type from DB
func (s EDPComponentService) GetEDPComponent(componentType string) (*query.EDPComponent, error) {
	log.Debug("start fetching EDP Component", zap.String("type", componentType))
	c, err := s.IEDPComponent.GetEDPComponent(componentType)
	if err != nil {
		if dberror.IsNotFound(err) {
			log.Debug("edp component wasn't found in DB", zap.String("name", componentType))
			return nil, nil
		}
		return nil, errors.Wrapf(err, "an error has occurred while fetching EDP Component by %v type from DB",
			componentType)
	}
	log.Info("edp component has been fetched from DB",
		zap.String("type", c.Type), zap.String("url", c.Url))
	return c, nil
}

//GetEDPComponents gets all EDP components from DB
func (s EDPComponentService) GetEDPComponents() ([]*query.EDPComponent, error) {
	log.Debug("start fetching EDP Components...")
	c, err := s.IEDPComponent.GetEDPComponents()
	if err != nil {
		return nil, errors.Wrap(err, "an error has occurred while fetching EDP Components from DB")
	}
	log.Info("edp components have been fetched", zap.Any("length", len(c)))

	for i, v := range c {
		modifyPlatformLinks(v.Url, c[i])
	}

	return c, nil
}

func modifyPlatformLinks(url string, c *query.EDPComponent) {
	replacer := strings.NewReplacer("{namespace}", context.Namespace)
	c.Url = fmt.Sprintf("%v%v", url, replacer.Replace(beego.AppConfig.String("projectMaskUrl")))
}
