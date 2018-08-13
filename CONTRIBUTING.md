# How to contribute

We are glad you are reading this, because we need volunteer developers to help
this project come to fruition.

If you don't have anything you are working on we have a list of newbie friendly
issues you can help out with.

If you haven't already, come find us on our mailing list. We want you working on
things you're excited about.

### Code of Conduct

Jenny, like most other open source projects, has a Code of Conduct that it
expects its contributors and core team members to adhere to. Our
[Code of Conduct](https://github.com/jennyservices/jenny/blob/master/CODE_OF_CONDUCT.md)
is available at the root of this repo. Any violation of the Code of Conduct
should be reported to
[Typeform Open-Source](https://open-source.typeform.com/to/xYAH1q) immidiately.

Here are some additional resources;

* Mailing List: [jenny-dev](https://groups.google.com/forum/#!forum/jenny-dev)
* Bug Tracker: [Issues](https://github.com/jennyservices/jenny/issues)
* Docs: [Docs](https://github.com/jennyservices/jenny/tree/master/docs)
* GoDoc:
  [![GoDoc](https://godoc.org/github.com/Typeform/jenny?status.svg)](https://godoc.org/github.com/Typeform/jenny)
* Travis-CI:
  [![Build Status](https://travis-ci.org/Typeform/jenny.svg?branch=master)](https://travis-ci.org/Typeform/jenny)
* Coverage:
  [![codecov](https://codecov.io/gh/Typeform/jenny/branch/master/graph/badge.svg)](https://codecov.io/gh/Typeform/jenny)

# How can I help?

## What should I know?

Please familiarize your self with why we built Jenny and some of the design
decisions behind it.

* [Introducing Jenny](https://ffbyt.es/YULg)

### Design Decisions

* [Generator](https://github.com/jennyservices/jenny/blob/master/docs/generator.md)

## Reporting Bugs

Jenny is in it's early stages of infancy, so there are bugs, we'd love your help
fixing them. In order to submit meaningful bug reports that we can take action
on

* Use a clear and descriptive title for the issue to identify the problem.
* Describe the exact steps which reproduce the problem in as many details as
  possible. Start by explaining what was your input service definition and your
  desired output language, what did you expect to happen what did?
* Provide specific examples to demonstrate the steps. Include links to files or
  GitHub projects, or copy/pasteable snippets, which you use in those examples.
  If you're providing snippets in the issue, use Markdown code blocks.
* Describe the behavior you observed after following the steps and point out
  what exactly is the problem with that behavior.
* Explain which behavior you expected to see instead and why.
* Include screenshots and animated GIFs which show you following the described
  steps and clearly demonstrate the problem.
* If the problem wasn't triggered by a specific action, describe what you were
  doing before the problem happened and share more information using the
  guidelines below.

### Suggesting Enhancements

Right now Jenny doesn't benefit from the same amount attention and engineering
resources as more popular projects like Kuberentes benefits from, therefore we
in this early stage we can't accept major enhancements.

## Your First Code Contribution

Unsure where to begin contributing to Atom? You can start by looking through
these `beginner` and `help-wanted` issues:

* [Beginner issues](https://github.com/jennyservices/jenny/labels/beginner) - issues
  which should only require a few lines of code, and a test or two.
* [Help wanted issues](https://github.com/jennyservices/jenny/labels/help%20wanted) -
  issues which should be a bit more involved than beginner issues.

## Your first Pull Request

We use github pull requests for accepting changes.

### Easiest way of getting something merged

We hope by now you have found/reported an bug you want to send a PR for.

First thing first, let's start by switching to a bug branch, let's assume you
are fixing issue #20, which is a bug about Jenny ignoring `format` fields in
definitions.

```
git checkout -b issue20
```

Second thing to do is adding `issue20_test.go`. Here you should add a reference
to the bug you are fixing. If one does not exist, describe the issue here.

```golang
// Bug report: https://github.com/jennyservices/jenny/issues/20
// Reported by: @marc-gr
// Issue: Now, if a field type is number, but its format is float, jenny ignores
// the format and set it as int. We should take into account the format when
// generating the code.
func TestIssue20(t *testing.T) {
  /* your tests here */
}
```

Run `go test ./...` make sure the tests fail, where you expect them to.

Commit your test that are failing.

```
golang: add tests for issue20

Proposed change adds tests for the golang package more specifically the
getSimpleType function which does not respect the format fields.

Reported-By: marc-gr
Signed-Off-By: sevki
```

Fix the bug! And run `go test ./...` until the tests are passing again.

When the code is fixed, commit your fix separately from your tests.

```
golang: fix getSimpleType to respect formats

When passed format and type getSimpleType function now respects the format field

Closes #20

Reported-By: marc-gr
Signed-Off-By: sevki
```

You can now submit your Pull Request! Our pull request has two commits, one adds
a new test case and the other one fixes it, the commit messages are clear and
states what they fix. A maintainer looking at this has all the information they
would need to merge this pull request.

## Style Guide

For all `Go` code, follow the
[Go style guide](https://github.com/golang/go/wiki/Style)

### Commit Message

Your commit message should follow this format.

```
{affectedPackage}: {action} {simple explanation}

{Work that was done in more detail}

{Closes #20}

R: {reviewers,}

Reported-By: {reporter}
Signed-Off-By: {fixer}
```

Actions: [add, fix, modify, remove]
