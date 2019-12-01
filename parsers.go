package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseSession(file *os.File) (*Session, error) {
	var sections []string
	scanner := bufio.NewScanner(file)

	scanner.Split(splitAt("\r\n\r\n"))
	for scanner.Scan() {
		sections = append(sections, scanner.Text())
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	if len(sections) != 4 {
		return nil, errors.New("invalid file, does not contain 4 sections")
	}

	hash := sha1.New()
	hash.Write([]byte(file.Name()))
	byteSum := hash.Sum(nil)
	sha1Hash := fmt.Sprintf("%x", byteSum)

	timeStamp, err := parseTimeStamp(file.Name())
	if err != nil {
		return nil, err
	}

	var generalSettings *GeneralSettings
	if len(sections) >= 4 {
		generalSettings, err = parseGeneralSettings(sections[3])
	}

	var weaponSettings *WeaponSettings
	if len(sections) >= 2 {
		weaponSettings, err = parseWeaponSettings(sections[1])
	}

	var sessionStats *Statistics
	if len(sections) >= 3 {
		sessionStats, err = parseSessionStats(sections[2])
	}

	var kills []*Kill
	if len(sections) >= 1 {
		kills, err = parseKills(sections[0])
	}

	session := &Session{
		SessionHash:     sha1Hash,
		Time:            timeStamp,
		GeneralSettings: generalSettings,
		WeaponSettings:  weaponSettings,
		Statistics:      sessionStats,
		Kills:           kills,
	}

	return session, nil
}

func parseGeneralSettings(s string) (*GeneralSettings, error) {
	var m map[string]string
	m = make(map[string]string)

	columns := strings.Split(s, "\r\n")
	for _, column := range columns {
		headingAndValue := strings.Split(column, ":,")
		if len(column) >= 1 {
			heading := headingAndValue[0]
			value := headingAndValue[1]
			m[heading] = value
		}
	}

	columnName := "Input Lag"
	inputLag, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Max FPS (config)"
	maxFPS, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Horiz Sens"
	horizSens, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Vert Sens"
	vertSens, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "FOV"
	fov, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Hide Gun"
	hideGun, err := strconv.ParseBool(m[columnName])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Crosshair Scale"
	crosshairScale, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	return &GeneralSettings{
		InputLag:       inputLag,
		MaxFPS:         maxFPS,
		SensScale:      m["Sens Scale"],
		HorizSens:      horizSens,
		VertSens:       vertSens,
		FOV:            fov,
		HideGun:        hideGun,
		Crosshair:      m["Crosshair"],
		CrosshairScale: crosshairScale,
		CrosshairColor: m["Crosshair Color"],
	}, nil
}

func parseWeaponSettings(s string) (settings *WeaponSettings, err error) {
	stringReader := strings.NewReader(s)
	csvReader := csv.NewReader(stringReader)
	csvReader.FieldsPerRecord = -1 // variable number of fields per record ...
	row := 0

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// don't do anything with the first row, which contains column names
		// TODO: Can CSVReader grab the column by name rather than by index?
		if row > 0 {
			settings = &WeaponSettings{}

			settings.Shots, err = strconv.ParseInt(record[1], 0, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Shots' column.", record[1]))
			}

			settings.Hits, err = strconv.ParseInt(record[2], 0, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as int64 for 'Hits' column.", record[2]))
			}

			settings.DamageDone, err = strconv.ParseFloat(record[3], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Damage Done' column.", record[3]))
			}

			settings.DamagePossible, err = strconv.ParseFloat(record[4], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Damage Possible' column.", record[4]))
			}

			// these fields aren't as important, and aren't on every scenario so we just default them if they error out.
			if len(record) >= 7 {
				settings.HorizSens, _ = strconv.ParseFloat(record[7], 0)
				settings.VertSens, _ = strconv.ParseFloat(record[8], 0)
				settings.FOV, _ = strconv.ParseFloat(record[9], 0)
				settings.HideGun, _ = strconv.ParseBool(record[10])
				settings.CrosshairScale, _ = strconv.ParseFloat(record[12], 0)
				settings.ADSSens, _ = strconv.ParseFloat(record[14], 0)
				settings.ADSZoomScale, _ = strconv.ParseFloat(record[15], 0)
			}
		}
		row = row + 1
	}

	return settings, nil
}

func parseSessionStats(s string) (*Statistics, error) {
	var m map[string]string
	m = make(map[string]string)

	columns := strings.Split(s, "\r\n")
	for _, column := range columns {
		headingAndValue := strings.Split(column, ":,")
		heading := headingAndValue[0]
		value := headingAndValue[1]
		m[heading] = value
	}

	columnName := "Kills"
	kills, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Deaths"
	deaths, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Fight Time"
	fightTime, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Avg TTK"
	avgTTK, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Damage Done"
	damageDone, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Damage Taken"
	damageTaken, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Midairs"
	midairs, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Midaired"
	midaired, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Directs"
	directs, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Directed"
	directed, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Distance Traveled"
	distanceTraveled, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	columnName = "Score"
	score, err := strconv.ParseFloat(m[columnName], 0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for '%s' column.", m[columnName], columnName))
	}

	statistics := &Statistics{
		Kills:            kills,
		Deaths:           deaths,
		FightTime:        fightTime,
		AvgTTK:           avgTTK,
		DamageDone:       damageDone,
		DamageTaken:      damageTaken,
		Midairs:          midairs,
		Midaired:         midaired,
		Directs:          directs,
		Directed:         directed,
		DistanceTraveled: distanceTraveled,
		Scenario:         m["Scenario"],
		Score:            score,
		Hash:             m["Hash"],
		GameVersion:      m["Game Version"],
	}
	return statistics, nil
}

func parseKills(s string) (kills []*Kill, err error) {
	stringReader := strings.NewReader(s)
	csvReader := csv.NewReader(stringReader)
	row := 0

	for {

		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// don't do anything with the first row, which contains column names
		// TODO: Can CSVReader grab the column by name rather than by index?
		if row > 0 {

			killNumber, err := strconv.ParseFloat(record[0], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Kill #' column.", record[0]))
			}

			shots, err := strconv.ParseFloat(record[5], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Shots' column.", record[5]))
			}

			hits, err := strconv.ParseFloat(record[6], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Hits' column.", record[6]))
			}

			accuracy, err := strconv.ParseFloat(record[7], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Accuracy' column.", record[7]))
			}

			damageDone, err := strconv.ParseFloat(record[8], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Damage Done' column.", record[8]))
			}

			damagePossible, err := strconv.ParseFloat(record[9], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Damage Possible' column.", record[9]))
			}

			efficiency, err := strconv.ParseFloat(record[10], 0)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as float for 'Efficiency' column.", record[10]))
			}

			cheated, err := strconv.ParseBool(record[11])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("can not parse '%s' as bool for 'Cheated' column.", record[11]))
			}

			var kill = &Kill{
				KillNumber:     killNumber,
				Timestamp:      record[1],
				Bot:            record[2],
				Weapon:         record[3],
				TTK:            record[4],
				Shots:          shots,
				Hits:           hits,
				Accuracy:       accuracy,
				DamageDone:     damageDone,
				DamagePossible: damagePossible,
				Efficiency:     efficiency,
				Cheated:        cheated,
			}

			kills = append(kills, kill)
		}

		row = row + 1
	}

	return kills, nil
}

func parseTimeStamp(filePath string) (time.Time, error) {
	// fun times with file names, the only place you can get the date...
	timeStampStartIndex := strings.LastIndex(filePath, " - ")
	timeStampStartIndex = timeStampStartIndex + 3
	timeStampString := filePath[timeStampStartIndex : timeStampStartIndex+19]
	timeStampString = strings.Replace(timeStampString, "-", " ", 1)

	dateTime := strings.Split(timeStampString, " ")
	if len(dateTime) != 2 {
		log.Fatalln("Time parsing did not result in two fields: date/time.")
	}

	date := dateTime[0]
	date = strings.Replace(date, ".", "-", -1)
	clock := dateTime[1]
	clock = strings.Replace(clock, ".", ":", -1)

	timeStampString = fmt.Sprintf("%s %s", date, clock)
	timestamp, err := time.Parse("2006-01-02 15:04:05", timeStampString)
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("unable to parse timestamp: " + timeStampString))
	}

	return timestamp, nil
}

func splitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), substring); i >= 0 {
			return i + len(substring), data[0:i], nil
		}

		// if at end of file with data, return the data
		if atEOF {
			return len(data), data, nil
		}

		return
	}
}
