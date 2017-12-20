/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

class Footer extends React.Component {
  render() {
    const currentYear = new Date().getFullYear();
    return (
      <footer className="nav-footer" id="footer">
        <section className="sitemap">
          <a href={this.props.config.baseUrl} className="nav-home">
            <img
              src={this.props.config.baseUrl + this.props.config.footerIcon}
              alt={this.props.config.title}
              width="66"
              height="58"
            />
          </a>
          <div>
            <h5>Docs</h5>
            <a href={this.props.config.baseUrl + 'docs/readme.html'}>
              Getting Started
            </a>
            <a href={this.props.config.baseUrl + 'docs/tutorials.html'}>
              Tutorials{' '}
            </a>
            <a href={'https://godoc.org/github.com/Typeform/jenny'}>Go Docs</a>
          </div>
          <div>
            <h5>Community</h5>
            <a
              href={
                this.props.config.baseUrl + this.props.language + '/users.html'
              }>
              User Showcase
            </a>

            <a
              href="https://groups.google.com/forum/#!forum/jenny-dev"
              target="_blank">
              Mailing List
            </a>
            <a href="https://github.com/typeform/jenny/issues" target="_blank">
              Issues
            </a>
          </div>
          <div>
            <h5>More</h5>
            <a href={this.props.config.baseUrl + 'blog'}>Blog</a>
            <a
              className="github-button"
              href={this.props.config.repoUrl}
              data-icon="octicon-star"
              data-count-href="/facebook/docusaurus/stargazers"
              data-show-count={true}
              data-count-aria-label="# stargazers on GitHub"
              aria-label="Star this project on GitHub">
              Star
            </a>
          </div>
        </section>

        <a
          href="https://typeform.github.io"
          target="_blank"
          className="fbOpenSource">
          Typeform Open Source
        </a>
        <section className="copyright">
          Copyright &copy; {currentYear} Typeform SL.
        </section>
      </footer>
    );
  }
}

module.exports = Footer;
