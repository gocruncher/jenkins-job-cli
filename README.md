#  Runs Jenkins job from the Command Line 
<meta name="google-site-verification" content="Wl2WZRolJ6omFNTQRguTy0GRQU41taSDq20n4Qgz05c" />

The utility starts a Jenkins build/job from the Command Line/Terminal.
An execution will be like this:

![terminal demo](assets/demo.gif)

## Install
Fetch the [latest release](https://github.com/gocruncher/jenkins-job-ctl/releases) for your platform:

#### Linux

```bash
sudo wget https://github.com/gocruncher/jenkins-job-ctl/releases/download/v1.0.1/jenkins-job-ctl-1.0.1-linux-amd64 -O /usr/local/bin/jj
sudo chmod +x /usr/local/bin/jj
```

#### OS X brew

```bash
brew tap gocruncher/tap
brew install jj
```
#### OS X bash
```bash
sudo curl -Lo /usr/local/bin/jj https://github.com/gocruncher/jenkins-job-ctl/releases/download/v1.0.1/jenkins-job-ctl-1.0.1-darwin-amd64
sudo chmod +x /usr/local/bin/jj
```

## Getting Started 

### Configure Access to Multiple Jenkins

```bash
jj set dev_jenkins --url "https://myjenkins.com" --login admin --token 11aa0926784999dab5  
```
where the token is available in your personal configuration page of the Jenkins. Go to the Jenkins Web Interface and click your name on the top right corner on every page, then click "Configure" to see your API token. 

In case, when Jenkins is available without authorization:
```bash
jj set dev_jenkins --url "https://myjenkins.com"  
```

or just run the following command in dialog execution mode:
```bash
jj set dev_jenkins
```


### Shell autocompletion

As a recommendation, you can enable shell autocompletion for convenient work. To do this, run following:
```bash
# for zsh completion:
echo 'source <(jj completion zsh)' >>~/.zshrc

# for bash completion:
echo 'source <(jj completion bash)' >>~/.bashrc
```
if this does not work for some reason, try following command that might help you to figure out what is wrong: 
```bash
jj completion check
```

### Examples
```bash
# Configure Access to the Jenkins
jj set dev-jenkins

# Start 'app-build' job in the current Jenkins
jj run app-build

# Start 'web-build' job in Jenkins named prod
jj run -n prod web-build

# makes a specific Jenkins name by default
jj use PROD  
```

## Futures
- cancellation job (Ctrl+C key)
- resize of the output (just press enter key)
- output of child jobs   

## Useful packages
- [cobra](https://github.com/spf13/cobra) - library for creating powerful modern CLI
- [chalk](https://github.com/chalk/chalk) – Terminal string styling done right
- [bar](https://github.com/superhawk610/bar) - Flexible ascii progress bar.

## Todos
- add authorization by login/pass and through the RSA key
- support of a terminal window resizing

## License
`jenkins-job-ctl` is open-sourced software licensed under the [MIT](LICENSE) license.
