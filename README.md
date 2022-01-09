<h1 align="center">Jenkins Notifier</h1>
<p align="center">Windows notification for jenkin builds</h2>

## Todos
- [x] Configuration file to store job information
- [x] Pause / unpause checking job status
- [ ] ✨GUI✨ (Ongoing)
- [ ] Check job's last status from GUI
- [ ] Edit configuration from GUI

## Configuration
<b>config.json</b>

|   Name   |   Type   |                           Description                           |
|----------|----------|-----------------------------------------------------------------|
| job      |  `array` | Array containing information about jobs                         |
| job.name | `string` | Name of the job <b>(should be unique)</b>                       |
| job.tag  | `string` | Not used                                                        |
| job.url  | `string` | Url of the job                                                  |
| interval | `number` | Interval between checks                                         |
| authType | `string` | Type of authentication.<br>Accepts either `password` or `token` |

<b>keys.json</b>

| Name     | Type     | Description                                                                                                                                                                                |
|----------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| username | `string` | Jenkins username                                                                                                                                                                           |
| apiKey   | `string` | If `authType` is set to `token`, then the api key is used to authenticate.<br>[How to get Jenkins API token](https://stackoverflow.com/a/45466184/15440719 "How to get Jenkins API token") |

## Building
```
go build
```
