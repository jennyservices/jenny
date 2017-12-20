/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

const CompLibrary = require('../../core/CompLibrary.js');
const MarkdownBlock = CompLibrary.MarkdownBlock; /* Used to read markdown */
const Container = CompLibrary.Container;
const GridBlock = CompLibrary.GridBlock;

const siteConfig = require(process.cwd() + '/siteConfig.js');

class Button extends React.Component {
  render() {
    return (
      <div className="pluginWrapper buttonWrapper">
        <a className="button" href={this.props.href} target={this.props.target}>
          {this.props.children}
        </a>
      </div>
    );
  }
}

Button.defaultProps = {
  target: '_self'
};

class HomeSplash extends React.Component {
  render() {
    return (
      <div className="homeContainer">
        <div className="homeSplashFade">
          <div className="wrapper homeWrapper">
            <div className="projectLogo">
              <img src={siteConfig.baseUrl + 'img/jenny.svg'} />
            </div>
            <div className="inner">
              <h2 className="projectTitle">
                {siteConfig.title}
                <small>{siteConfig.tagline}</small>
              </h2>
              <div className="section promoSection">
                <div className="promoRow">
                  <div className="pluginRowBlock">
                    <Button href={siteConfig.baseUrl + 'docs' + '/readme.html'}>
                      Get Started
                    </Button>
                    <Button href={'https://github.com/Typeform/jenny'}>
                      GitHub
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

class Index extends React.Component {
  render() {
    let language = this.props.language || 'en';
    const showcase = siteConfig.users
      .filter(user => {
        return user.pinned;
      })
      .map(user => {
        return (
          <a href={user.infoLink}>
            <img src={user.image} title={user.caption} />
          </a>
        );
      });

    return (
      <div>
        <HomeSplash language={language} />
        <div className="mainContainer">
          <Container padding={['bottom', 'top']}>
            <GridBlock
              align="center"
              contents={[
                {
                  content:
                    'jenny will generate your code for you! <br> You can extend it to generate code for any language you want!',
                  image: siteConfig.baseUrl + 'img/code.svg',
                  imageAlign: 'top',
                  title: 'Code Generator'
                },
                {
                  content:
                    'jenny makes testing your application and playing around with it a breeze. <br>Find out how?',
                  image: siteConfig.baseUrl + 'img/bug.svg',
                  imageAlign: 'top',
                  title: 'Debugger'
                }
              ]}
              layout="fourColumn"
            />
          </Container>

          <div
            className="productShowcaseSection paddingBottom"
            style={{ textAlign: 'center' }}>
            <h2>Production Ready!</h2>
            <MarkdownBlock>
              Built on top of [go-kit](http://gokit.io/) jenny services are
              production ready!
            </MarkdownBlock>
          </div>

          {/*          <Container padding={['bottom', 'top']} background="light">
            <GridBlock
              contents={[
                {
                  content: 'Talk about learning how to use this',
                  image: siteConfig.baseUrl + 'img/jenny-code-gen.png',
                  imageAlign: 'right',
                  title: 'Learn How'
                }
              ]}
            />
          </Container>

          <Container padding={['bottom', 'top']} id="try">
            <GridBlock
              contents={[
                {
                  content:
                    "Jenny takes a service definition and generates code for you, currently it supports generating server-side Go code with JS on the way. You can write your own code generators to extend jenny's code generation capabilities.",
                  image: siteConfig.baseUrl + 'img/jenny-code-gen.png',
                  imageAlign: 'left',
                  title: 'Generate Code'
                }
              ]}
            />
          </Container>

          <Container padding={['bottom', 'top']} background="dark">
            <GridBlock
              contents={[
                {
                  content:
                    'This is another description of how this project is useful',
                  image: siteConfig.baseUrl + 'img/docusaurus.svg',
                  imageAlign: 'right',
                  title: 'Debug'
                }
              ]}
            />
          </Container>*/}

          <div className="productShowcaseSection paddingBottom">
            <h2>{"Who's Using This?"}</h2>
            <p>This project is used by all these people</p>
            <div className="logos">{showcase}</div>
            <div className="more-users">
              <a
                className="button"
                href={
                  siteConfig.baseUrl + this.props.language + '/' + 'users.html'
                }>
                More {siteConfig.title} Users
              </a>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

module.exports = Index;
