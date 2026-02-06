package requester

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const workerCount int64 = 10
const sequenceCount int64 = 500

const requestJitterRandNVal = 11 // also in time.Millisecond
const requestJitterMin = 100 * time.Millisecond

const targetIdHubUrl = "http://localhost:20100/api/identity"
const targetDidUrl = "http://localhost:20000"
const targetIdHubApiKey = "YWRtaW4.adminKey"

type Requester struct {
	Cursor *RequesterCursor
	Logger *slog.Logger
	Client *http.Client
}

type RequesterRun struct {
	Did                        string
	ParticipantRequest         ParticipantRequest
	ParticipantRequestResponse any
	ParticipantRequestFailed   bool

	DidDocumentResponse any
	DidDocumentFailed   bool

	KeyPairResponse       any
	KeyPairResponseFailed bool
}

func NewRequester(defaultLogger *slog.Logger) *Requester {
	return &Requester{
		Cursor: NewRequesterCursor(),
		Logger: defaultLogger,
		Client: &http.Client{
			Timeout: time.Millisecond * 500,
			Transport: &http.Transport{
				// Connection pooling
				MaxIdleConns:        1000,
				MaxIdleConnsPerHost: 200,
				MaxConnsPerHost:     0, // unlimited (safe for local tests)

				// Keep connections around
				IdleConnTimeout: 90 * time.Second,

				// Timeouts to avoid goroutine leaks
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,

				// Faster local dialing
				DialContext: (&net.Dialer{
					Timeout:   2 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,

				// Disable compression if CPU is a concern in tests
				DisableCompression: true,
			},
		},
	}
}

func (requester *Requester) RunAll() {
	var wg sync.WaitGroup
	for workerId := range workerCount {
		workerLogger := requester.Logger.With(
			slog.Int64("workerId", workerId),
		)

		wg.Go(func() {
			requester.RunSequence(workerId, sequenceCount, workerLogger)
		})
	}
	wg.Wait()
}

func (requester *Requester) RunSequence(workerId int64, count int64, workerLogger *slog.Logger) {
	workerLogger.Info("Starting Sequence", "Count", count)
	sequenceRunData := make([]*RequesterRun, 0, count)
	for runId := range count {
		if runId%10 == 0 {
			workerLogger.Info(
				"Progress",
				slog.Int64("run", runId),
				slog.Int64("total", count),
			)
		}
		resultRun := requester.RunSingle(workerId, runId)
		sequenceRunData = append(sequenceRunData, resultRun)
	}
	workerLogger.Info("Finished Sequence")
}

func (requester *Requester) RunSingle(workerId int64, runId int64) *RequesterRun {
	jitter := time.Duration(rand.Intn(requestJitterRandNVal)) * time.Millisecond
	time.Sleep(requestJitterMin + jitter)
	didString := requester.Cursor.GetNextDid()

	logger := requester.Logger.With(
		slog.Int64("workerId", workerId),
		slog.Int64("runId", runId),
		slog.String("did", didString),
	)

	reqRun := RequesterRun{
		Did:                        didString,
		ParticipantRequest:         NewParticipantRequest(didString),
		ParticipantRequestResponse: "",
		ParticipantRequestFailed:   false,
		DidDocumentResponse:        "",
		DidDocumentFailed:          false,
		KeyPairResponse:            "",
		KeyPairResponseFailed:      false,
	}

	if err := requester.SendRequest(&reqRun); err != nil {
		reqRun.ParticipantRequestFailed = true
		logger.Warn("participant request failed", "error", err)
		return &reqRun
	}

	time.Sleep(requestJitterMin + jitter)

	if err := requester.RetrieveDidDocument(&reqRun); err != nil {
		reqRun.DidDocumentFailed = true
		logger.Warn("requesting did document failed", "error", err)
		return &reqRun
	}

	time.Sleep(requestJitterMin + jitter)

	err := requester.CheckDidVerificationMethods(&reqRun)

	if err == nil {
		if rand.Intn(100) > 98 {
			// also log some good runs
			if err := requester.RetrieveKeyPairForRun(&reqRun); err != nil {
				logger.Warn("could not retrieve keypair for run", "error", err)
				reqRun.KeyPairResponseFailed = true
			}
			requester.WriteRunToFile(&reqRun, "ok", logger)
		}

		// all as expected, verificationMethod is not empty.
		return &reqRun
	}

	reqRun.DidDocumentFailed = true
	logger.Warn("checking for verification methods failed! this might be a case.", "error", err)

	// when we get here the verification method was not there.
	// so try to gather the keypair and write all of it to a log file

	if err := requester.RetrieveKeyPairForRun(&reqRun); err != nil {
		logger.Warn("could not retrieve keypair for run", "error", err)
		reqRun.KeyPairResponseFailed = true
	}
	requester.WriteRunToFile(&reqRun, "failed", logger)

	return &reqRun
}

func (requester *Requester) SendRequest(run *RequesterRun) error {
	body, err := json.Marshal(run.ParticipantRequest)
	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest(
		http.MethodPost,
		targetIdHubUrl+"/v1alpha/participants",
		bodyReader,
	)

	req.Header.Add("X-Api-Key", targetIdHubApiKey)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return fmt.Errorf("failed to create request for did %v: %w", run.Did, err)
	}

	res, err := requester.Client.Do(req)
	if err != nil {
		return fmt.Errorf("requesting participant failed for did %v: %w", run.Did, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-ok status %v", res.StatusCode)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse body for did %v: %w", run.Did, err)
	}

	var responseJson any
	if err := json.Unmarshal(responseBody, &responseJson); err != nil {
		return fmt.Errorf("failed to parse body to json for did %v: %w", run.Did, err)
	}

	run.ParticipantRequestResponse = responseJson
	return nil
}

func (requester *Requester) RetrieveDidDocument(run *RequesterRun) error {
	didParts := strings.Split(run.Did, ":")
	didSubString := didParts[len(didParts)-1]

	req, err := http.NewRequest(
		http.MethodGet,
		targetDidUrl+"/"+didSubString,
		nil,
	)

	req.Host = "consumer-idhub"

	if err != nil {
		return fmt.Errorf("failed to create request for did document: %v", err)
	}

	res, err := requester.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request did document: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-ok status %v", res.StatusCode)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("body response read error %v", err)
	}

	var responseJson any
	if err := json.Unmarshal(responseBody, &responseJson); err != nil {
		return fmt.Errorf("failed to parse did document body to json for did %v: %w", run.Did, err)
	}

	run.DidDocumentResponse = responseJson
	return nil
}

func (requester *Requester) CheckDidVerificationMethods(run *RequesterRun) error {
	didDocument, ok := run.DidDocumentResponse.(map[string]any)
	if !ok {
		return fmt.Errorf("assertion failed: didDocumentResponse is not map[string]any")
	}

	verificationMethods, hasVerMethods := didDocument["verificationMethod"]

	if !hasVerMethods {
		return fmt.Errorf("no verification method map entry")
	}

	methodList, ok := verificationMethods.([]any)

	if !ok {
		return fmt.Errorf("assertion of list of for verificationMethod failed")
	}

	if len(methodList) == 0 {
		return fmt.Errorf("list of verificationMethod is empty")
	}

	return nil
}

func (requester *Requester) RetrieveKeyPairForRun(run *RequesterRun) error {
	encodedDid := base64.URLEncoding.EncodeToString([]byte(run.Did))

	req, err := http.NewRequest(
		http.MethodGet,
		targetIdHubUrl+"/v1alpha/participants/"+encodedDid+"/keypairs",
		nil,
	)

	req.Header.Add("X-Api-Key", targetIdHubApiKey)

	if err != nil {
		return fmt.Errorf("failed to create request for keypair for did %v: %w", run.Did, err)
	}

	res, err := requester.Client.Do(req)
	if err != nil {
		return fmt.Errorf("requesting keypair failed for did %v: %w", run.Did, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-ok status %v", res.StatusCode)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse body for did %v: %w", run.Did, err)
	}

	var responseJson any
	if err := json.Unmarshal(responseBody, &responseJson); err != nil {
		return fmt.Errorf("failed to parse key pair body to json for did %v: %w", run.Did, err)
	}

	run.KeyPairResponse = responseJson
	return nil
}

func (requester *Requester) WriteRunToFile(run *RequesterRun, suffix string, logger *slog.Logger) {
	if len(suffix) == 0 {
		logger.Warn("empty suffix not allowed for logs")
		return
	}

	logFileName := "logs/" + suffix + "/" + time.Now().Format(time.RFC3339) + "_" + run.Did + "_" + suffix + ".json"

	if err := os.MkdirAll("logs/"+suffix, 0777); err != nil {
		logger.Warn("failed to create suffix logs folder", "error", err)
		return
	}

	fileContent, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		logger.Warn("failed to marshal json of req run", "error", err)
		return
	}

	os.WriteFile(logFileName, fileContent, 0644)
}

type RequesterCursor struct {
	Count    int
	CountMut *sync.RWMutex
}

func NewRequesterCursor() *RequesterCursor {
	return &RequesterCursor{
		Count:    0,
		CountMut: &sync.RWMutex{},
	}
}

func (cursor *RequesterCursor) GetNextDid() string {
	cursor.CountMut.Lock()
	defer cursor.CountMut.Unlock()

	var sb strings.Builder
	sb.WriteString("did:web:consumer-idhub:test-id-")
	cursor.Count++
	nextInt := cursor.Count
	sb.WriteString(strconv.FormatInt(int64(nextInt), 10))

	return sb.String()
}
