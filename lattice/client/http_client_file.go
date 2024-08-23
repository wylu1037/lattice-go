package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/wylu1037/lattice-go/common/types"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (api *httpApi) UploadFile(_ context.Context, chainId, filePath string) (*types.UploadFileResponse, error) {
	log.Debug().Msgf("开始上传文件到链上，chainId: %s, filePath: %s", chainId, filePath)
	uploadPath := fmt.Sprintf("%s/%s", api.GinServerUrl, "beforeSign")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}(file)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadPath, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set(headerContentType, writer.FormDataContentType()) // fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())
	req.Header.Set(headerChainID, chainId)
	if api.jwtApi != nil {
		token, _ := api.jwtApi.GetToken()
		req.Header.Set(headerAuthorize, fmt.Sprintf("Bearer %s", token))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	uploadFileResponse := new(types.UploadFileResponse)
	if err := json.Unmarshal(responseData, uploadFileResponse); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response body")
		return nil, err
	}
	log.Debug().Msgf("结束上传文件【%s】到链上", filePath)
	return uploadFileResponse, nil
}

func (api *httpApi) DownloadFile(_ context.Context, cid, filePath string) error {
	log.Debug().Msgf("开始从链上下载文件【%s】", cid)
	downloadUrl := fmt.Sprintf("%s/download?cid=%s", api.GinServerUrl, cid)

	downloadReq, err := http.NewRequest(http.MethodGet, downloadUrl, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request for download")
		return fmt.Errorf("failed to create request: %w", err)
	}
	if api.jwtApi != nil {
		token, _ := api.jwtApi.GetToken()
		downloadReq.Header.Set(headerAuthorize, fmt.Sprintf("Bearer %s", token))
		downloadReq.Header.Set(headerContentType, "multipart/form-data; charset=UTF-8")
		downloadReq.Header.Set(headerConnection, "close")
	}

	client := &http.Client{}
	resp, err := client.Do(downloadReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to download file")
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close response body")

		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载文件【%s】失败，Http状态吗为: %d", cid, resp.StatusCode)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(outFile)
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to copy file data")
	}

	log.Debug().Msgf("结束从链上下载文件【%s】", cid)
	return nil
}
