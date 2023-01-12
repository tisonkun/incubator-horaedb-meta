// Copyright 2022 CeresDB Project Authors. Licensed under Apache-2.0.

package http

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CeresDB/ceresmeta/pkg/log"
	"github.com/CeresDB/ceresmeta/server/member"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ForwardClient struct {
	member *member.Member
	client *http.Client
	port   int
}

func NewForwardClient(member *member.Member, port int) *ForwardClient {
	return &ForwardClient{
		member: member,
		client: getForwardedHTTPClient(),
		port:   port,
	}
}

func getForwardedHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func (s *ForwardClient) getForwardedAddr(ctx context.Context) (string, bool, error) {
	member, err := s.member.GetLeader(ctx)
	if err != nil {
		return "", false, errors.WithMessage(err, "get forwarded addr")
	}
	if member.IsLocal {
		return "", true, nil
	}
	httpAddr, err := formatHTTPAddr(member.Leader.Endpoint, s.port)
	if err != nil {
		return "", false, errors.WithMessage(err, "format http addr")
	}
	return httpAddr, false, nil
}

func (s *ForwardClient) forwardToLeader(req *http.Request) (*http.Response, bool, error) {
	addr, isLeader, err := s.getForwardedAddr(req.Context())
	if err != nil {
		log.Error("get forward addr failed", zap.Error(err))
		return nil, false, err
	}
	if isLeader {
		return nil, true, nil
	}

	// Update remote host
	req.RequestURI = ""
	if req.TLS == nil {
		req.URL.Scheme = "http"
	} else {
		req.URL.Scheme = "https"
	}
	req.URL.Host = addr

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error("forward client send request failed", zap.Error(err))
		return nil, false, err
	}

	return resp, false, nil
}

// formatHttpAddr convert grpcAddr(http://127.0.0.1:8831) httpPort(5000) to httpAddr(127.0.0.1:5000).
func formatHTTPAddr(grpcAddr string, httpPort int) (string, error) {
	leaderAddr := strings.Split(grpcAddr, ":")
	if len(leaderAddr) != 3 {
		return "", errors.WithMessagef(ErrParseLeaderAddr, "gprc addr:%s", grpcAddr)
	}
	leaderAddr[2] = strconv.Itoa(httpPort)
	leaderAddr = append(leaderAddr[:0], leaderAddr[0:]...)
	httpAddr := strings.Join(leaderAddr, ":")
	return httpAddr, nil
}