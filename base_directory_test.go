// (c) 2014 John R. Lenton. See LICENSE.

package xdg

import (
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestXDGd(t *testing.T) { TestingT(t) }

type xdgdSuite struct {
	home string
	env1 string
	val1 string
	env2 string
	val2 string
	dir  *XDGDir
}

var _ = Suite(&xdgdSuite{})

func (s *xdgdSuite) SetUpTest(c *C) {
	s.home = os.Getenv("HOME")
	s.env1 = "go_xdg_one"
	s.env2 = "go_xdg_two"
	s.val1 = "something"
	s.val2 = "one:two:three"
	s.dir = &XDGDir{s.env1, s.val1, s.env2, s.val2}
}

func (s *xdgdSuite) TestHomePrefersEnviron(c *C) {
	err := os.Setenv(s.env1, "algo")
	c.Assert(err, IsNil)
	defer os.Setenv(s.env1, "")
	h := s.dir.Home()
	c.Check(h, Equals, "algo")
}

func (s *xdgdSuite) TestHomeUsesDefault(c *C) {
	h := s.dir.Home()
	c.Check(h, Matches, s.home+".*"+s.val1)
}

func (s *xdgdSuite) TestDirsPrefersEnviron(c *C) {
	err := os.Setenv(s.env1, "cero")
	c.Assert(err, IsNil)
	defer os.Setenv(s.env1, "")
	err = os.Setenv(s.env2, "uno:dos")
	c.Assert(err, IsNil)
	defer os.Setenv(s.env2, "")
	hs := s.dir.Dirs()
	c.Check(hs, DeepEquals, []string{"cero", "uno", "dos"})
}

func (s *xdgdSuite) TestDirsSkipsEmpty(c *C) {
	err := os.Setenv(s.env2, "::")
	c.Assert(err, IsNil)
	defer os.Setenv(s.env2, "")
	hs := s.dir.Dirs()
	c.Check(hs, HasLen, 1)
}

func (s *xdgdSuite) TestDirsUsesDefault(c *C) {
	hs := s.dir.Dirs()
	c.Assert(hs, HasLen, 4)
	c.Check(hs[1:], DeepEquals, strings.Split(s.val2, ":"))
	c.Check(hs[0], Matches, s.home+".*"+s.val1)
}

// now repeat all the tests, but without the HOME environ.
type xdgdNoHomeSuite struct {
	xdgdSuite
}

var _ = Suite(&xdgdNoHomeSuite{})

func (s *xdgdNoHomeSuite) SetUpTest(c *C) {
	s.xdgdSuite.SetUpTest(c)
	os.Setenv("HOME", "")
}

func (s *xdgdNoHomeSuite) TearDownTest(c *C) {
	os.Setenv("HOME", s.home)
}

// and for these tests, an entirely fake HOME
type xdgdFHSuite struct {
	xdgdSuite
	real_home string
}

var _ = Suite(&xdgdFHSuite{})

func (s *xdgdFHSuite) SetUpTest(c *C) {
	s.real_home = os.Getenv("HOME")
	home := c.MkDir()
	os.Setenv("HOME", home)
	s.xdgdSuite.SetUpTest(c)
	s.val2 = c.MkDir() + ":" + c.MkDir() + ":" + c.MkDir()
	s.dir = &XDGDir{s.env1, s.val1, s.env2, s.val2}
}

func (s *xdgdFHSuite) TearDownTest(c *C) {
	os.Setenv("HOME", s.real_home)
}

func (s *xdgdFHSuite) TestFind(c *C) {
	vs := strings.Split(s.val2, ":")
	res1 := "stuff"
	exp1 := filepath.Join(s.home, s.val1, res1)
	res2 := "things/that"
	exp2 := filepath.Join(vs[1], res2)
	res3 := "more"
	exp3 := filepath.Join(vs[2], res3)
	for _, d := range []string{exp1, exp2, exp3} {
		err := os.MkdirAll(d, 0700)
		c.Assert(err, IsNil, Commentf(d))
	}
	for _, it := range []struct {
		res string
		exp string
	}{{res1, exp1}, {res2, exp2}, {res3, exp3}} {
		rv, err := s.dir.Find(it.res)
		c.Assert(err, IsNil)
		c.Check(rv, Equals, it.exp)
	}
	_, err := s.dir.Find("missing")
	c.Check(err, NotNil)
}

func (s *xdgdFHSuite) TestEnsureFirst(c *C) {
	// creates it if missing
	rv1, err := s.dir.Ensure("missing/file")
	c.Assert(err, IsNil)
	_, err = os.Stat(rv1)
	c.Check(err, IsNil)
	c.Check(rv1, Matches, s.home+".*"+"missing/file")
	// just gets it if existing
	rv2, err := s.dir.Ensure("missing/file")
	c.Assert(err, IsNil)
	c.Check(rv2, Equals, rv1)
}

func (s *xdgdFHSuite) TestEnsureFirstFailures(c *C) {
	_, err := s.dir.Ensure(strings.Repeat("*", 1<<9) + "/" + strings.Repeat("*", 1<<9))
	c.Assert(err, NotNil)
}
