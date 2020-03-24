c:\tools\pscp -i %USERPROFILE%\.ssh\pi.ppk -r -v dist pi@192.168.1.158:/home/pi/hires-vector
c:\tools\pscp -i %USERPROFILE%\.ssh\pi.ppk hires-vector.go pi@192.168.1.158:/home/pi/hires-vector
rem list files:
rem c:\tools\pscp -i %USERPROFILE%\.ssh\pi.ppk -ls pi@sidcloud.net:/home/pi