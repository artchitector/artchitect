package infrastructure

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

/*
стандартные параметры генерации, отлаженные на RTX3060 12Gb
Artchitect1 settings: 640x960 (x4) -> 2560x3840

можно рисовать и картинки 3328 x 5120 (это 832x1280 до апскейла). Рисует 50 секунд. Смысла нет.
И вроде рисует хуже в больших разрешениях
*/
const (
	TxtFilename  = "list.txt"
	ImageWidth   = 640 // должно быть кратно 64
	ImageHeight  = 960 // должно быть кратно 64 и быть х1.5 от ширины
	ImageSteps   = 50
	ImageUpscale = 4
)

// AI - обёртка над Invoke.AI (который обёртка над Stable Diffusion), дающая создание картинок по prompt
type AI struct {
	invokeAIPath     string
	pathFinderRegexp *regexp.Regexp
}

func NewAI(invokeAIPath string) *AI {
	// регулярное выражение для поиска файла в выводе от InvokeAI выглядит так
	// .*(\/home\/user\/invoke-ai\/invokeai_v2.3.0\/outputs\/[0-9\.]+png).*
	regexpPath := strings.ReplaceAll(invokeAIPath, "/", "\\/")
	fullRegexp := fmt.Sprintf(".*(%s\\/outputs\\/[0-9\\.]+png).*", regexpPath)
	return &AI{invokeAIPath: invokeAIPath, pathFinderRegexp: regexp.MustCompile(fullRegexp)}
}

func (ai *AI) GenerateImage(ctx context.Context, seed uint, prompt string) ([]byte, error) {
	if err := ai.prepareListTxt(seed, prompt); err != nil {
		return nil, errors.Wrap(err, "[ai] ОШИБКА СОЗДАНИЯ TXT-ФАЙЛА")
	}

	filename, err := ai.execute(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[ai] ОШИБКА ВЫЗОВА ИИ")
	} else if filename == "" {
		log.Debug().Msgf("[ai] РИСОВАНИЕ ПРЕРВАНО")
		return nil, errors.New("[ai] ОСТАНОВ")
	}
	log.Debug().Msgf("[ai] НАЙДЕН ИСКОМЫЙ ФАЙЛ КАРТИНЫ - %s", filename)
	imgData, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "[ai] ОШИБКА ЧТЕНИЯ ФАЙЛА %s", filename)
	}
	log.Debug().Msgf("[ai] ФАЙЛ ПРОЧИТАН. РАЗМЕР %dКБ", len(imgData)/1000)
	return imgData, err
}

// prepareListTxt - Invoke.AI будет читать данные из файла list.txt
// На этом шаге надо заполнить файл list.txt данными для нового изображения
func (ai *AI) prepareListTxt(seed uint, prompt string) error {
	/*
		строка файла выглядит так:
		{prompt} -S{seed} -W{width} -H{height} -s{steps} -U{upscale}
		пример:
		dark side of the moon,human,similar -S3780028127 -W640 -H960 -s50 -U4
	*/
	line := strings.Join(
		[]string{
			prompt,
			fmt.Sprintf("-S%d", seed),
			fmt.Sprintf("-W%d", ImageWidth),
			fmt.Sprintf("-H%d", ImageHeight),
			fmt.Sprintf("-s%d", ImageSteps),
			fmt.Sprintf("-U%d", ImageUpscale),
		},
		" ",
	)
	listTxtPath := path.Join(ai.invokeAIPath, TxtFilename)
	err := os.WriteFile(listTxtPath, []byte(line), 0644)
	return err
}

// prepareCmd - сборка консольной команды, которая запустит Invoke.AI
func (ai *AI) prepareCmd() string {

	return strings.Join(
		[]string{
			fmt.Sprintf("export INVOKEAI_ROOT=%s && ", ai.invokeAIPath),
			fmt.Sprintf("%s/.venv/bin/python", ai.invokeAIPath),
			fmt.Sprintf("%s/.venv/bin/invoke.py", ai.invokeAIPath),
			//fmt.Sprintf(`--from_file "%s/list.txt"`, ai.invokeAIPath),
		},
		" ",
	)
}

func (ai *AI) execute(ctx context.Context) (string, error) {
	// https://www.yellowduck.be/posts/reading-command-output-line-by-line

	/*
		https://github.com/invoke-ai/InvokeAI
		invokeAIPath - пусть до локальной установки Invoke.AI

		Вся команда выглядит так (берём python из InvokeAI и им запускаем скрипт):
		INVOKEAI_ROOT=<invokeAIPath> \
			<invokeAIPath>/.venv/bin/python \
			<invokeAIPath>/.venv/bin/invoke.py \
			--from_file <invokeAIPath>/list.txt
	*/
	cmd := exec.Command(
		fmt.Sprintf("%s/.venv/bin/python", ai.invokeAIPath),
		fmt.Sprintf("%s/.venv/bin/invoke.py", ai.invokeAIPath),
		fmt.Sprintf(`--from_file`),
		fmt.Sprintf(`%s/list.txt`, ai.invokeAIPath),
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("INVOKEAI_ROOT=%s", ai.invokeAIPath))

	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = cmd.Stdout
	scanner := bufio.NewScanner(r)
	done := make(chan struct{})
	var filename string

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			var found bool
			found, filename = ai.checkLineAndGetFile(line)
			if found {
				break
			}
		}
		done <- struct{}{}
	}()

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		log.Info().Msgf("[ai] УБИТЬ ПРОЦЕСС %d", cmd.Process.Pid)
		if err := cmd.Process.Kill(); err != nil {
			return "", errors.Wrapf(err, "[ai] ПРОЦЕСС %d НЕ УБИТ", cmd.Process.Pid)
		} else {
			log.Debug().Msgf("[ai] ПРОЦЕСС %d УБИТ", cmd.Process.Pid)
			return "", nil
		}
	case <-done:
		break
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	if filename == "" {
		return filename, errors.New("[ai] НЕ НАЙДЕНА СТРОКА С ФАЙЛОМ")
	}

	return filename, nil
}

func (ai *AI) checkLineAndGetFile(line string) (found bool, filename string) {
	log.Debug().Msgf("[ai] СТРОКА: %s", line)
	match := ai.pathFinderRegexp.FindStringSubmatch(line)
	if len(match) > 1 {
		return true, match[1]
	}
	return false, ""
}
