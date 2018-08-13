/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

const CompLibrary = require('../../core/CompLibrary.js');
const Container = CompLibrary.Container;
const GridBlock = CompLibrary.GridBlock;

const siteConfig = require(process.cwd() + '/siteConfig.js');

class Help extends React.Component {
  render() {
    const supportLinks = [
      {
        content:
          'Learn more using the [documentation on this site.](/jenny/docs/readme.html)',
        title: 'Browse Docs'
      },
      {
        content:
          'Ask questions about the jenny on our [mailing list.](https://groups.google.com/forum/#!forum/jenny-dev)',
        title: 'Join the community'
      },
      {
        content:
          "Find out what's new with jenny on our [blog](http://engineering.typeform.com/jenny/blog/)",
        title: 'Stay up to date'
      }
    ];

    return (
      <div className="docMainWrapper wrapper">
        <Container className="mainContainer documentContainer postContainer">
          <div className="post">
            <header className="postHeader">
              <h2>Want to help?</h2>
            </header>
            <p>
              Please visit our{' '}
              <a href="https://github.com/jennyservices/jenny/blob/master/CONTRIBUTING.md">
                contribution guidelines
              </a>{' '}
              to see how you can start helping.
            </p>
            <p>
              Issues
              <ul>
                <li>
                  <a href="https://github.com/jennyservices/jenny/issues?q=is%3Aissue+is%3Aopen+label%3Abeginner">
                    Beginner issues
                  </a>
                </li>
              </ul>
            </p>
            <p>This project is maintained by a dedicated group of people;</p>
            <ul>
              <li>
                <a href="https://github.com/sevki">@sevki</a>
              </li>
            </ul>
            <GridBlock contents={supportLinks} layout="threeColumn" />
          </div>
        </Container>
      </div>
    );
  }
}

module.exports = Help;
