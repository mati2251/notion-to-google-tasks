# notion-to-google-tasks
## Purpose
Notion to Google Tasks is a console application designed to simplify the process of transferring tasks and lists from Notion to Google Tasks 
## Usage
Run subcomand config to configure the application
```bash
notion-to-google-tasks config
```
After that command create config file in directory `~/.config/notion-to-google-tasks/config.yml`. The example of config file is in test.sample.yaml <br>
If you want sync your databases with list just subcommand sync
```bash
notion-to-google-tasks sync
```
## Tech stack
- [Go](https://golang.org/)
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [Notion API](https://developers.notion.com/)
- [Google Tasks API](https://developers.google.com/tasks)