package app

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	fsnotify "github.com/fsnotify/fsnotify"
	"github.com/vzemtsov/config-reloader/config"
	"github.com/vzemtsov/config-reloader/pkg/metrics"
)

func Run(cfg *config.Config) error {

	err := checks(cfg)
	if err != nil {
		return err
	}

	if len(cfg.VolumeDirs) > 0 {
		volumeDirWatcher(cfg)
	}
	if len(cfg.VolumeDirsArchive) > 0 {
		for _, vda := range cfg.VolumeDirsArchive {
			unarchiveDir(vda, cfg)
		}
		if *cfg.InitMode {
			log.Println("Init mode completed")
			return nil
		}
		volumeDirArchiveWatcher(cfg)
	}

	return nil

}

func volumeDirWatcher(cfg *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if !isValidEvent(event) {
					continue
				}

				log.Println("ConfigMap or Secret updated")
				sendWebHook(cfg)
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				metrics.WatcherErrors.Inc()
				log.Println("Error:", err)
			}
		}
	}()

	for _, d := range cfg.VolumeDirs {
		log.Printf("Watching directory: %q", d)
		err = watcher.Add(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func volumeDirArchiveWatcher(cfg *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if !isValidEvent(event) {
					continue
				}

				err := unarchiveFile(event.Name, cfg)
				if err != nil {
					log.Println("Error:", err)
				}

				log.Println("ConfigMap or Secret updated")
				sendWebHook(cfg)
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				metrics.WatcherErrors.Inc()
				log.Println("Error:", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	for _, d := range cfg.VolumeDirsArchive {
		log.Printf("Watching directory (with unarchive): %q", d)
		err = watcher.Add(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func checks(cfg *config.Config) error {
	if (len(cfg.VolumeDirs) < 1) && (len(cfg.VolumeDirsArchive) < 1) {
		return fmt.Errorf("%s", "Missing volume-dir or volume-dir-archive")
	}

	if len(cfg.Webhook.Urls) < 1 {
		return fmt.Errorf("%s", "Missing webhook-url")
	}

	if *cfg.InitMode && (len(cfg.VolumeDirsArchive) < 1) {
		return fmt.Errorf("%s", "init-mode work only with volume-dir-archive")
	}

	if *cfg.InitMode && (len(cfg.VolumeDirs) > 0) {
		return fmt.Errorf("%s", "init-mode don't work with volume-dir")
	}

	return nil
}

func isValidEvent(event fsnotify.Event) bool {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return false
	}

	if filepath.Base(event.Name) != "..data" {
		return false
	}
	return true
}

func sendWebHook(cfg *config.Config) {
	for _, h := range cfg.Webhook.Urls {
		begun := time.Now()
		req, err := http.NewRequest(*cfg.Webhook.Method, h.String(), nil)
		if err != nil {
			metrics.SetFailureMetrics(h.String(), "client_request_create")
			log.Println("Error:", err)
			continue
		}
		userInfo := h.User
		if userInfo != nil {
			if password, passwordSet := userInfo.Password(); passwordSet {
				req.SetBasicAuth(userInfo.Username(), password)
			}
		}

		successfulReloadWebhook := false

		for retries := *cfg.Webhook.Retries; retries != 0; retries-- {
			log.Printf("Performing webhook request (%d/%d)", retries, *cfg.Webhook.Retries)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				metrics.SetFailureMetrics(h.String(), "client_request_do")
				log.Println("Error:", err)
				time.Sleep(time.Second * 10)
				continue
			}
			resp.Body.Close()
			metrics.RequestsByStatusCode.WithLabelValues(h.String(), strconv.Itoa(resp.StatusCode)).Inc()
			if resp.StatusCode != *cfg.Webhook.StatusCode {
				metrics.SetFailureMetrics(h.String(), "client_response")
				log.Println("error:", "Received response code", resp.StatusCode, ", expected", cfg.Webhook.StatusCode)
				time.Sleep(time.Second * 10)
				continue
			}

			metrics.SetSuccessMetrics(h.String(), begun)
			log.Println("successfully triggered reload")
			successfulReloadWebhook = true
			break
		}

		if !successfulReloadWebhook {
			metrics.SetFailureMetrics(h.String(), "retries_exhausted")
			log.Println("error:", "Webhook reload retries exhausted")
		}
	}

}

func unarchiveDir(path string, cfg *config.Config) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullFilePath := path + "/" + file.Name()
		err := unarchiveFile(fullFilePath, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func unarchiveFile(path string, cfg *config.Config) error {

	if path[len(path)-3:] != ".gz" {
		return fmt.Errorf("File %s is not a .gz archive. Do nothing", path)
	}

	gzipFile, err := os.Open(path)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Uncompress to a writer. We'll use a file writer
	outFileName := *cfg.DirForUnarchive + "/" + filepath.Base(path)[0:len(filepath.Base(path))-3]
	// fmt.Println(outFileName)

	outfileWriter, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer outfileWriter.Close()

	// Copy contents of gzipped file to output file
	_, err = io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return err
	}
	return nil
}
