# Udemy Kurs: Go: Schnelle & sichere Webanwendungen programmieren

## PostgreSQL - Soda / Pop / Fizz

Achtung: Das Beispiel von database.yml aus [https://gobuffalo.io/documentation/database/configuration/](https://gobuffalo.io/documentation/database/configuration/) ist für die Verwendung von "nur" Soda nicht korrekt. Dieses Beispiel ist für die Verwendung im Gesamtframework Go Buffalo. Das Problem bei dem Beispiel ist, dass es die Templatesyntax von Go Buffalo verwendet ({{ und }}), welche für YAML aber nicht korrekt ist. Daher muss man die Teile test: und production: entfernen bzw. auskommentieren, sonst kann soda das database.yml nicht korrekt interpretieren und beschwert sich, dass `There is no connection named development defined!`

Hier ein "funktionierendes" database.yml für die Verwendung mit nur Soda:

```yaml
development:
  dialect: postgres
  database: mygowebapp
  user: <<<enter-your-postgres-user-here (normally postgres)>>>
  password: <<<enter-your-password-here>>>
  host: 127.0.0.1
  pool: 5
```

## MailHog

Der verwendete Mailserver in dem Kurs ist MailHog. Man installiert es und muss es dann auch immer davor aufrufen (im Startmenü Mailhog eingeben, dann sollte der Shortcut dafür auftauchen). Der Mailserver lauscht auf Port 1025, auf Port 8025 kommt man zur Übersicht der Mails.

## Admincredentials

"Verbotenerweise" ist das Kennwort "password", die Mailadresse kann man aus der Tabelle users in der Postgres-Datenbank herausfinden.
