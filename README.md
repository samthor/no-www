Go binary which should be run on "www." hostnames.
This does only a few things:

* removes all "www." prefix from the domain
* redirects to https:// only
* only responds to GET requests

For example:

  http://www.example.com => https://example.com
  http://www.www.foo.bar => https://foo.bar
  https://test.com => 404
  https://www.test.com/hello => https://www.test.com/hello
