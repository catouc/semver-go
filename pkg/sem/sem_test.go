package sem

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestVersionKindIsValid(t *testing.T) {
	testData := []struct {
		name   string
		inKind Kind
		err    error
	}{
		{
			name:   "Major",
			inKind: 1,
			err:    nil,
		},
		{
			name:   "Minor",
			inKind: 2,
			err:    nil,
		},
		{
			name:   "Patch",
			inKind: 3,
			err:    nil,
		},
		{
			name:   "TooHigh",
			inKind: 10,
			err:    ErrInvalidKind,
		},
		{
			name:   "TooLow",
			inKind: -1,
			err:    ErrInvalidKind,
		},
		{
			name:   "Invalid",
			inKind: 0,
			err:    ErrInvalidKind,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			err := td.inKind.IsValid()
			if err != td.err {
				t.Errorf("expected %s got %s", td.err, err)
			}
		})
	}
}

func TestParseKind(t *testing.T) {
	testData := []struct {
		name string
		in   string
		out  Kind
		err  error
	}{
		{
			name: "ValidKind",
			in:   "major",
			out:  Major,
			err:  nil,
		},
		{
			name: "ValidKind",
			in:   "minor",
			out:  Minor,
			err:  nil,
		},
		{
			name: "ValidKind",
			in:   "patch",
			out:  Patch,
			err:  nil,
		},
		{
			name: "ValidKind",
			in:   "invalid",
			out:  Invalid,
			err:  ErrInvalidKind,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			kind, err := ParseKind(td.in)
			if err != td.err {
				t.Errorf("expected error %s got %s", td.err, err)
			}
			if kind != td.out {
				t.Errorf("expected kind %d got %d", td.out, kind)
			}
		})
	}
}

func TestGetNextVersion(t *testing.T) {
	testData := []struct {
		name             string
		inCurrentVersion *Ver
		inKind           Kind
		outNextVersion   *Ver
		err              error
	}{
		{
			name:             "Major",
			inCurrentVersion: &Ver{"v", 1, 0, 0},
			inKind:           Major,
			outNextVersion:   &Ver{"v", 2, 0, 0},
			err:              nil,
		},
		{
			name:             "Minor",
			inCurrentVersion: &Ver{"v", 1, 0, 0},
			inKind:           Minor,
			outNextVersion:   &Ver{"v", 1, 1, 0},
			err:              nil,
		},
		{
			name:             "Patch",
			inCurrentVersion: &Ver{"v", 1, 0, 1},
			inKind:           Patch,
			outNextVersion:   &Ver{"v", 1, 0, 2},
			err:              nil,
		},
		{
			name:             "MajorResetMinorPatch",
			inCurrentVersion: &Ver{"v", 1, 1, 1},
			inKind:           Major,
			outNextVersion:   &Ver{"v", 2, 0, 0},
			err:              nil,
		},
		{
			name:             "MinorResetPatch",
			inCurrentVersion: &Ver{"v", 1, 1, 1},
			inKind:           Minor,
			outNextVersion:   &Ver{"v", 1, 2, 0},
			err:              nil,
		},
		{
			name:             "NotAVersion",
			inCurrentVersion: &Ver{"v", 1, 0, 0},
			inKind:           -1,
			outNextVersion:   &Ver{"v", 1, 0, 0},
			err:              ErrInvalidKind,
		},
		{
			name:             "NotAVersion2",
			inCurrentVersion: &Ver{"v", 1, 0, 0},
			inKind:           10,
			outNextVersion:   &Ver{"v", 1, 0, 0},
			err:              ErrInvalidKind,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			err := td.inCurrentVersion.Next(td.inKind)
			if err != td.err {
				t.Errorf("expected error %s got %s", td.err, err)
			}
			if !reflect.DeepEqual(*td.inCurrentVersion, *td.outNextVersion) {
				t.Errorf("expected version %+v got %+v", *td.outNextVersion, *td.inCurrentVersion)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	testData := []struct {
		name       string
		inVersion  string
		outVersion *Ver
		err        error
	}{
		{
			name:       "ValidVersionMajor",
			inVersion:  "1.0.0",
			outVersion: &Ver{"", 1, 0, 0},
			err:        nil,
		},
		{
			name:       "ValidVersionMinor",
			inVersion:  "0.1.0",
			outVersion: &Ver{"", 0, 1, 0},
			err:        nil,
		},
		{
			name:       "ValidVersionHighNumbers",
			inVersion:  "9999999.999999.99999",
			outVersion: &Ver{"", 9999999, 999999, 99999},
			err:        nil,
		},
		{
			name:       "InvalidVersionNegativeNumbers",
			inVersion:  "-1.-1.-1",
			outVersion: nil,
			err:        ErrParseVersionFault,
		},
		{
			name:       "ValidVersionPrefix",
			inVersion:  "v1.0.0",
			outVersion: &Ver{"v", 1, 0, 0},
			err:        nil,
		},
		{
			name:       "InvalidVersionPrefixEmoji",
			inVersion:  "☹️1.0.0",
			outVersion: nil,
			err:        ErrParseVersionFault,
		},
		{
			name:       "InvalidVersionNaNMajor",
			inVersion:  "a.b.c",
			outVersion: nil,
			err:        ErrParseVersionFault,
		},
		{
			name:       "InvalidVersionNaNMinor",
			inVersion:  "1.b.c",
			outVersion: nil,
			err:        ErrParseVersionFault,
		},
		{
			name:       "InvalidVersionNaNPatch",
			inVersion:  "1.0.c",
			outVersion: nil,
			err:        ErrParseVersionFault,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			ver, err := ParseVersion(td.inVersion)
			if err != td.err {
				t.Errorf("expected error %s got %s", td.err, err)
			}
			if ver != nil {
				if !reflect.DeepEqual(*ver, *td.outVersion) {
					t.Errorf("expected version struct %+v got %+v", *td.outVersion, *ver)
				}
			}
		})
	}
}

func TestGetAllVersions(t *testing.T) {
	testData := []struct {
		name                string
		cmdList             []*exec.Cmd
		ignoreNonSemVerTags bool
		out                 []*Ver
		err                 error
	}{
		{
			name:    "ValidVersionListSingleTag",
			cmdList: []*exec.Cmd{exec.Command("git", "tag", "v1.0.0")},
			out:     []*Ver{{"v", 1, 0, 0}},
			err:     nil,
		},
		{
			name:    "ValidVersionListSingleTagNoPrefix",
			cmdList: []*exec.Cmd{exec.Command("git", "tag", "1.0.0")},
			out:     []*Ver{{"", 1, 0, 0}},
			err:     nil,
		},
		{
			name: "ValidVersionListMultiTag",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v1.0.0"),
				exec.Command("git", "tag", "v1.1.0"),
				exec.Command("git", "tag", "v2.0.0"),
				exec.Command("git", "tag", "v2.0.1"),
				exec.Command("git", "tag", "v2.1.0"),
			},
			out: []*Ver{
				{"v", 1, 0, 0},
				{"v", 1, 1, 0},
				{"v", 2, 0, 0},
				{"v", 2, 0, 1},
				{"v", 2, 1, 0},
			},
			err: nil,
		},
		{
			name: "ValidVersionListMultiTagNoPrefix",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "1.0.0"),
				exec.Command("git", "tag", "1.1.0"),
				exec.Command("git", "tag", "2.0.0"),
				exec.Command("git", "tag", "2.0.1"),
				exec.Command("git", "tag", "2.1.0"),
			},
			out: []*Ver{
				{"", 1, 0, 0},
				{"", 1, 1, 0},
				{"", 2, 0, 0},
				{"", 2, 0, 1},
				{"", 2, 1, 0},
			},
			err: nil,
		},
		{
			name: "ValidVersionListReversed",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "2.1.0"),
				exec.Command("git", "tag", "2.0.1"),
				exec.Command("git", "tag", "2.0.0"),
				exec.Command("git", "tag", "1.1.0"),
				exec.Command("git", "tag", "1.0.0"),
			},
			out: []*Ver{
				{"", 1, 0, 0},
				{"", 1, 1, 0},
				{"", 2, 0, 0},
				{"", 2, 0, 1},
				{"", 2, 1, 0},
			},
			err: nil,
		},
		{
			name: "ValidVersionListMultiPrefix",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v2.1.0"),
				exec.Command("git", "tag", "a2.0.1"),
				exec.Command("git", "tag", "h2.0.0"),
				exec.Command("git", "tag", "d1.1.0"),
				exec.Command("git", "tag", "w1.0.0"),
			},
			out: []*Ver{
				{"w", 1, 0, 0},
				{"d", 1, 1, 0},
				{"h", 2, 0, 0},
				{"a", 2, 0, 1},
				{"v", 2, 1, 0},
			},
			err: nil,
		},
		{
			name: "ValidAndInvalidVersionListMultiPrefixIgnoreInvalidTags",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v2.1.0"),
				exec.Command("git", "tag", "a2.0.1"),
				exec.Command("git", "tag", "a.2.0.1"),
				exec.Command("git", "tag", "h2.0.0"),
				exec.Command("git", "tag", "d1.1.0"),
				exec.Command("git", "tag", "a-.2.0.1"),
				exec.Command("git", "tag", "w1.0.0"),
			},
			ignoreNonSemVerTags: true,
			out: []*Ver{
				{"w", 1, 0, 0},
				{"d", 1, 1, 0},
				{"h", 2, 0, 0},
				{"a", 2, 0, 1},
				{"v", 2, 1, 0},
			},
			err: nil,
		},
		{
			name: "InvalidVersionListNegativeVersions",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v-1.-1.-1")},
			out: nil,
			err: ErrParseVersionFault,
		},
		{
			name: "InvalidVersionListEmoji",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "☹️1.1.1")},
			out: nil,
			err: ErrParseVersionFault,
		},
		{
			name: "InvalidVersionListEmojiIgnoreInvalidTags",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "☹️1.1.1")},
			ignoreNonSemVerTags: true,
			out:                 []*Ver{},
			err:                 nil,
		},
		{
			name:    "NoVersions",
			cmdList: []*exec.Cmd{},
			out:     []*Ver{},
			err:     nil,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := setupTemporaryGitRepository(tmpDir, td.cmdList); err != nil {
				t.Fatal(err)
			}
			versionList, err := GetAllVersions(tmpDir, td.ignoreNonSemVerTags)
			if err != td.err {
				t.Errorf("expected error %s got %s", td.err, err)
			}
			if versionList != nil {
				if len(versionList) != len(td.out) {
					t.Errorf("version list size expected %d got %d", len(versionList), len(td.out))
				}
				for i, version := range versionList {
					if !reflect.DeepEqual(*version, *td.out[i]) {
						t.Errorf("expected version %+v got %+v", *td.out[i], *version)
					}
				}
			}
		})
	}
}

func TestGetLatestVersion(t *testing.T) {
	testData := []struct {
		name                string
		cmdList             []*exec.Cmd
		ignoreNonSemVerTags bool
		out                 *Ver
		err                 error
	}{
		{
			name:    "ValidVersionListSingleTag",
			cmdList: []*exec.Cmd{exec.Command("git", "tag", "v1.0.0")},
			out:     &Ver{"v", 1, 0, 0},
			err:     nil,
		},
		{
			name:    "ValidVersionListSingleTagNoPrefix",
			cmdList: []*exec.Cmd{exec.Command("git", "tag", "1.0.0")},
			out:     &Ver{"", 1, 0, 0},
			err:     nil,
		},
		{
			name:    "NoVersions",
			cmdList: []*exec.Cmd{},
			out:     nil,
			err:     ErrNoVersionsAvailable,
		},
		{
			name: "ValidVersionListMultiTag",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v1.0.0"),
				exec.Command("git", "tag", "v1.1.0"),
				exec.Command("git", "tag", "v2.0.0"),
				exec.Command("git", "tag", "v2.0.1"),
				exec.Command("git", "tag", "v2.1.0"),
			},
			out: &Ver{"v", 2, 1, 0},
			err: nil,
		},
		{
			name: "ValidVersionListMultiTagNoPrefix",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "1.0.0"),
				exec.Command("git", "tag", "1.1.0"),
				exec.Command("git", "tag", "2.0.0"),
				exec.Command("git", "tag", "2.0.1"),
				exec.Command("git", "tag", "2.1.0"),
			},
			out: &Ver{"", 2, 1, 0},
			err: nil,
		},
		{
			name: "ValidVersionListReversed",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "2.1.0"),
				exec.Command("git", "tag", "2.0.1"),
				exec.Command("git", "tag", "2.0.0"),
				exec.Command("git", "tag", "1.1.0"),
				exec.Command("git", "tag", "1.0.0"),
			},
			out: &Ver{"", 2, 1, 0},
			err: nil,
		},
		{
			name: "ValidVersionListMultiPrefix",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v2.1.0"),
				exec.Command("git", "tag", "a2.0.1"),
				exec.Command("git", "tag", "h2.0.0"),
				exec.Command("git", "tag", "d1.1.0"),
				exec.Command("git", "tag", "w1.0.0"),
			},
			out: &Ver{"v", 2, 1, 0},
			err: nil,
		},
		{
			name: "ValidAndInvalidVersionListMultiPrefixIgnoreInvalidTags",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v2.1.0"),
				exec.Command("git", "tag", "a2.0.1"),
				exec.Command("git", "tag", "a.2.0.1"),
				exec.Command("git", "tag", "h2.0.0"),
				exec.Command("git", "tag", "d1.1.0"),
				exec.Command("git", "tag", "a-.2.0.1"),
				exec.Command("git", "tag", "w1.0.0"),
			},
			ignoreNonSemVerTags: true,
			out:                 &Ver{"v", 2, 1, 0},
			err:                 nil,
		},
		{
			name: "InvalidVersionListNegativeVersions",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "v-1.-1.-1")},
			out: nil,
			err: ErrParseVersionFault,
		},
		{
			name: "InvalidVersionListEmoji",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "☹️1.1.1")},
			out: nil,
			err: ErrParseVersionFault,
		},
		{
			name: "InvalidVersionListEmojiIgnoreInvalidTags",
			cmdList: []*exec.Cmd{
				exec.Command("git", "tag", "☹️1.1.1")},
			ignoreNonSemVerTags: true,
			out:                 &Ver{},
			err:                 ErrNoVersionsAvailable,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := setupTemporaryGitRepository(tmpDir, td.cmdList); err != nil {
				t.Fatal(err)
			}
			version, err := GetLatestVersion(tmpDir, td.ignoreNonSemVerTags)
			if err != td.err {
				t.Errorf("expected error %s got %s", td.err, err)
			}
			if version != nil {
				if !reflect.DeepEqual(*version, *td.out) {
					t.Errorf("expected version %+v got %+v", *td.out, *version)
				}
			}
		})
	}
}

func TestGetAllVersionsInvalidDir(t *testing.T) {
	_, err := GetAllVersions(filepath.Join(".", "NaD"), false)
	if errors.Is(err, os.ErrInvalid) {
		t.Errorf("expected err to be os.PathError got: %s", err)
	}
}

func TestGetAllVersionsNoGitDir(t *testing.T) {
	_, err := GetAllVersions(t.TempDir(), false)
	if err == nil {
		t.Error("expected err got nil")
	}
}

// we're testing Sprintf for coverage so this is useless ¯\_(ツ)_/¯
func TestVerString(t *testing.T) {
	ver := Ver{"v", 1, 0, 0}
	if ver.String() != "v1.0.0" {
		t.Errorf("expected version string to be v1.0.0 got: %s", ver.String())
	}
}

func setupTemporaryGitRepository(dir string, additionalCommands []*exec.Cmd) error {
	initCommands := []*exec.Cmd{
		exec.Command("git", "init"),
		exec.Command("git", "config", "user.email", "test@example.com"),
		exec.Command("git", "config", "user.name", "test"),
		exec.Command("git", "commit", "--allow-empty", "-m", "\"Initial Commit\""),
	}
	commandList := append(initCommands, additionalCommands...)
	for _, cmd := range commandList {
		cmd.Dir = dir
		errBuf := bytes.Buffer{}
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed setup: %s: %s", err, errBuf.String())
		}
	}
	return nil
}
