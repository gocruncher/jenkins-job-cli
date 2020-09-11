# Jbuilder

Jbuilder is a simple command-line utility which just runs
any Jenkins job, view will be like this:

![terminal demo](assets/demo.gif)

## Installation
 
```
brew tap gocruncher/tap
brew install jb
```

## Quick start 

### Configure Access to Multiple Jenkins.

```
jb set dev_jenkins --url "https://myjenkins.com" --login admin --token 11aa0926784999dab5  
```
where the token is available in your personal configuration page of the Jenkins. Click your name on the top right corner on every page, then click "Configure" to see your API token. (The URL $root/me/configure is a good shortcut.) You can also change your API token from here.
In case, when Jenkins is available without authorization:
```
jb set dev_jenkins --url "https://myjenkins.com"  
```

Or just run the following command in dialog execution mode:
```
jb set dev_jenkins
```


### Shell autocompletion

you can enable shell autocompletion for convenient work. To do this, run following:
```
# for zsh completion:
echo 'source <(jb completion zsh)' >>~/.zshrc

# for bash completion:
echo 'source <(jb completion bash)' >>~/.bashrc
```
if this does not work for some reason, try following command that might help you to figure out what is wrong: 
```
jb completion check
```

### Usage
```
# run backend-app job in the current Jenkis 
jb run backend-app  

# run frontend job in the PROD Jenkins
jb -n PROD run frontend

# makes a specific Jenkins name by default
jb use PROD  
```

## Futures
- cancellation job (^C key)
- resize of the output (just press enter key)
- output of child jobs   

## Useful packages
- [cobra](https://github.com/spf13/cobra) - library for creating powerful modern CLI
- [chalk](https://github.com/chalk/chalk) – Terminal string styling done right
- [bar](https://github.com/superhawk610/bar) - Flexible ascii progress bar.

## Todos
- add authorization by login/pass and through the RSA key
- support of a terminal window resizing