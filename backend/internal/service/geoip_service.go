package service

import (
	"errors"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIP interface {
	Lookup(ip string) (country string, city string, err error)
}

type GeoIPService struct {
	db *geoip2.Reader
}

var (
	ErrGeoDBNotWorking = errors.New("geodb is not correctly initialized")
)

func NewGeoIPService(
	db *geoip2.Reader,
) *GeoIPService {
	return &GeoIPService{
		db: db,
	}
}

func (s *GeoIPService) Lookup(
	ip string,
) (string, string, error) {
	parsedIP := net.ParseIP(ip)
	log.Print(parsedIP)
	if parsedIP == nil {
		return "", "", errors.New("invalid IP address")
	}

	if s == nil || s.db == nil {
		return "", "", ErrGeoDBNotWorking
	}

	record, err := s.db.City(parsedIP)
	if err != nil {
		return "", "", err
	}
	log.Print(record)

	country := record.Country.IsoCode
	city := record.City.Names["en"]

	return country, city, nil
}
