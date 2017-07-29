// Released under an MIT license. See LICENSE.

// Generated by generate.oh

package boot

var Script string = `
define _connect_: syntax (conduit-name name) e = {
	define conduit-maker: e::eval conduit-name
	syntax (left right) e = {
		define p: conduit-maker
		define ec-ex-chan: channel 1
		spawn {
			define ec-ex: e::eval: _rew_: quasiquote: block {
				public (unquote name) = (unquote p)
				unquote left
			}
			p::_writer_close_
			ec-ex-chan::write ec-ex
		}
		block {
			define right-ec-ex: e::eval: _rew_: quasiquote: block {
				public _stdin_ = (unquote p)
				unquote right
			}
			p::_reader_close_
			define left-ec-ex: (ec-ex-chan::readlist)::head
			define ex: or (left-ec-ex::tail) (right-ec-ex::tail)
			if ex: throw ex
			return: and (left-ec-ex::head) (right-ec-ex::head)
		}
	}
}
define _redirect_: syntax (name mode closer) = {
	syntax (c cmd) e = {
		define c: e::eval c
		define f = ()
		if (not: or (is-channel c) (is-pipe c)) {
			set mode: e::eval mode
			if (and (eq mode: quote w) (exists -i c)) {
				throw: exception "${c} exists"
			}
			set f: open mode c
			set c = f
		}
		define ec-ex: e::eval: _rew_: quasiquote: block {
			public (unquote name) (unquote c)
			unquote cmd
		}
		if (not: is-null f): (f::_get_ closer)
		define ex: ec-ex::tail
		if ex: throw ex
		return: ec-ex::head
	}
}
define ...: method (: args) = {
	cd _origin_
	define path: args::head
	if (eq 2: args::length) {
		cd path
		set path: args::get 1
	}
	while true {
		define abs: symbol: "/"::join $PWD path
		if (exists abs): return abs
		if (eq $PWD /): return path
		cd ..
	}
}
define and: syntax (: lst) e = {
	define r = false
	while (not: is-null lst) {
		set r: e::eval: lst::head
		if (not r): return r
		set lst: lst::tail
	}
	return r
}
define _append_stderr_: _redirect_ _stderr_ "a" _writer_close_
define _append_stdout_: _redirect_ _stdout_ "a" _writer_close_
define _backtick_: syntax (cmd) e = {
	define p: pipe
	define ec-ex-chan: channel 1
	spawn {
		define ec-ex: e::eval: _rew_: quasiquote: block {
			public _stdout_ = (unquote p)
			unquote cmd
		}
		p::_writer_close_
		ec-ex-chan::write ec-ex
	}
	define r: cons () ()
	define c = r
	while (define l: p::readline) {
		c::set-tail: cons l ()
		set c: c::tail
	}
	p::_reader_close_
	define ex: ((ec-ex-chan::readlist)::head)::tail
	if ex: throw ex
	return: r::tail
}
define catch: syntax (name: clause) e = {
	define args: list name (quote throw)
	define body: list (quote throw) name
	if (is-null clause) {
		set body: list body
	} else {
		set body: clause::append body
	}
	define handler: e::eval {
		list (quote method) args (quote =) @body
	}
	define _return: e::eval (quote return)
	define _throw = throw
	e::public throw: method (condition) = {
		public throw = _throw
		_return: handler condition _throw
	}
}
define _channel_stderr_: _connect_ channel _stderr_
define _channel_stdout_: _connect_ channel _stdout_
define clobber: builtin (: args) = {
	tee @args >/dev/null
}
define coalesce: syntax (: lst) e = {
	while (and (not: is-null: lst::tail) (not: resolves: lst::head)) {
		set lst: lst::tail
	}
	return: e::eval: lst::head
}
define echo: builtin (: args) = {
	if (is-null args) {
		_stdout_::write: symbol ""
	} else {
		_stdout_::write @(for args symbol)
	}
}
define error: builtin (: args) =: _stderr_::write @args
define for: method (l m) = {
	define r: cons () ()
	define c = r
	while (not: is-null l) {
		c::set-tail: cons (m: l::head) ()
		set c: c::tail
		set l: l::tail
	}
	return: r::tail
}
define glob: builtin (: args) =: return args
define import-sem: channel 1
define import: method (module-path) = {
	catch unused {
		import-sem::readlist
	}
	import-sem::write ()
	define import-return = return
	define module-name = module-path
	define module: method (name) = {
		if (_root_::_modules_::has name) {
			define module-object: _root_::_modules_::_get_ name
			import-sem::readlist
			import-return module-object
		}
		set module-name = name
	}
	define module-object: object {
		source module-path
	}
	_root_::_modules_::_set_ module-name module-object
	import-sem::readlist
	import-return module-object
}
define is-list: method (l) = {
	if (is-null l): return false
	if (not: is-cons l): return false
	if (is-null: l::tail): return true
	is-list: l::tail
}
define is-text: method (t) =: or (is-string t) (is-symbol t)
# Last File Manager cd.
define lcd: method () e = {
	catch unused {
		return
	}
	if (not: which lfm >/dev/null) {
		return
	}
	public $LFMPATHFILE = "/tmp/lfm-${_pid_}.path"
	lfm -1
	if (not: exists $LFMPATHFILE) {
		return
	}
	define f: open "r" $LFMPATHFILE
	define dn: f::readline
	f::close
	rm $LFMPATHFILE
	e::eval: quasiquote: cd (unquote dn)
}
# Change working dir to last dir on exit from lf.
define lf: builtin (:args) e = {
	catch unused {
		return
	}
	if (not: which lf >/dev/null) {
		return
	}
	define temp-file @(_backtick_: mktemp)
	command lf "--last-dir-path=${temp-file}" @args
	if (not: test -f temp-file) {
		return
	}
	define dn @(_backtick_: cat temp-file)
	rm temp-file
	if dn {
		e::eval: quasiquote: cd (unquote dn)
	}
}
define module: method (name) = # no-op
define map: syntax (: literal) e = {
    define o: object
    for literal: method (entry) = {
        define k: entry::head
        if (is-list k) {
            set k: e::eval k
	}
        define v: entry::get 1
        if (eq 1: v::length) {
            set v: v::head
        }
        o::_set_ k: e::eval v
    }
    public del = _del_
    public get = _get_
    public set = _set_
    return o
}
# TODO: Replace with builtin rather than invoking bc.
define math: method (S) e = {
	catch ex {
		set ex::type = "error/syntax"
		set ex::message = "Malformed expression: ${S}"
	}

	float @(_backtick_ (block {
		echo "scale=6"
		write: symbol S
	} | bc))
}
define object: syntax (: body) e = {
	e::eval: cons (quote block): body::append: quote: context
}
define or: syntax (: lst) e = {
	define r = false
	while (not: is-null lst) {
		set r: e::eval: lst::head
		if r: return r
		set lst: lst::tail
	}
	return r
}
define _pipe_stderr_: _connect_ pipe _stderr_
define _pipe_stdout_: _connect_ pipe _stdout_
define printf: method (f: args) =: echo: (string f)::sprintf @args
define _process_substitution_: syntax (:args) e = {
	define chans = ()
	define fifos = ()
	define cmd: for args: method (arg) = {
		if (not: is-cons arg): return arg
		if (eq (quote _substitute_stdin_) (arg::head)) {
			define chan: channel 1
			define fifo: temp-fifo
			spawn {
				chan::write: e::eval: _rew_: quasiquote {
					_redirect_stdin_ {
						unquote fifo
						unquote: arg::tail
					}
				}
			}
			set chans: cons chan chans
			set fifos: cons fifo fifos
			return fifo
		}
		if (eq (quote _substitute_stdout_) (arg::head)) {
			define chan: channel 1
			define fifo: temp-fifo
			spawn {
				chan::write: e::eval: _rew_: quasiquote {
					_redirect_stdout_ {
						unquote fifo
						unquote: arg::tail
					}
				}
			}
			set chans: cons chan chans
			set fifos: cons fifo fifos
			return fifo
		}
		return arg
	}
	define main-ec-ex: e::eval: _rew_ cmd
	define ec-exs ()
	for chans: method (chan) = {
		set ec-exs: cons ((chan::readlist)::head) ec-exs
	}
	set ec-exs: cons main-ec-ex ec-exs
	for ec-exs: method (ec-ex) = {
		define ex: ec-ex::tail
		if ex: throw ex
	}
	rm @fifos
	return: main-ec-ex::head
}
define quasiquote: syntax (cell) e = {
	if (not: is-cons cell): return cell
	if (is-null cell): return cell
	if (eq (quote unquote): cell::head): return: e::eval: cell::get 1
	cons {
		e::eval: list (quote quasiquote) (cell::head)
		e::eval: list (quote quasiquote) (cell::tail)
	}
}
define readline: builtin () =: _stdin_::readline
define readlist: builtin () =: _stdin_::readlist
define _redirect_stderr_: _redirect_ _stderr_ "w" _writer_close_
define _redirect_stdin_: _redirect_ _stdin_ "r" _reader_close_
define _redirect_stdout_: _redirect_ _stdout_ "w" _writer_close_
define _rew_: method (command) = {    # return exception wrapper => rew
	quasiquote: (method () = {
		define ec: status false
		catch ex {
			return: cons ec ex
		}
		set ec: unquote command
		throw ()
	})
}
define source: syntax (name) e = {
	define basename: e::eval name
	define paths = ()
	define name = basename

	if (has $OHPATH): set paths: ":"::split $OHPATH
	while (and (not: is-null paths) (not: exists name)) {
		set name: "/"::join (paths::head) basename
		set paths: paths::tail
	}

	if (not: exists name): set name = basename

	define f: open r- name

	define r: cons () ()
	define c = r
	while (define l: f::readlist) {
		c::set-tail: cons l ()
		set c: c::tail
	}
	set c: r::tail
	f::close

	define rval: status 0
	define eval-list: method (first rest) o = {
		if (is-null first): return rval
		set rval: e::eval first
		eval-list (rest::head) (rest::tail)
	}
	eval-list (c::head) (c::tail)
	return rval
}
define umask: method (: args) = {
	catch ex {
		if (eq "_umask_: command not found" ex::message) {
			set ex::message = "not implemented on ${_platform_}"
		}
	}

	if (is-null args) {
		define omask: _umask_
		printf "%03o" omask
		return omask
	}
	_umask_ @args
}
define write: method (: args) =: _stdout_::write @args
public quote: syntax (cell) =: return cell
_root_::public _modules_: object
_sys_::public complete: method self (first completing) = {
	return ()
}
_sys_::public _exception: method (t s m l f) e = {
	object {
		public type = t
		public status = s
		public message = m
		public line = l
		public file = f
	}
}
# The generator method exception can be called in three ways:
# exception <message>
# exception <value> <message>
# exception <type> <value> <message>
# If not provided value defaults to status false and type to error/runtime.
_sys_::public exception: method (message :args) e = {
	define s: status false
	define t: symbol "error/runtime"
	if (not: is-null args) {
		set s: args::head
		set args: args::tail
	}
	if (not: is-null args) {
		set t = s
		set s: args::head
		set args: args::tail
	}
	_exception t s message {
		public line: e::eval: get-line-number
		public file: e::eval: get-source-file
	}
}
_sys_::public get-completions: method self (first completing) = {
	catch unused {
		return ()
	}
	self::complete first completing
}
_sys_::public get-prompt: method self (suffix) = {
	catch unused {
		return suffix
	}
	self::prompt suffix
}
_sys_::public prompt: method (suffix) = {
	define dirs: "/"::split $PWD
	return: ""::join (dirs::get: integer "-1") suffix
}
_sys_::public throw: method (c) = {
	error: ": "::join c::file c::line c::type c::message
	fatal c::status
}

exists ("/"::join $HOME .ohrc) && source ("/"::join $HOME .ohrc)

`

//go:generate ./generate.oh
//go:generate go fmt generated.go
