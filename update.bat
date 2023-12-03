@echo off
:wait
tasklist /fi "imagename eq boofutils.exe" 2>nul | find /i /n "boofutils.exe">nul
if "%!e(string=boofutils_new.exe)rrorlevel%!"(string=B:\git\Boofutils\boofutils.exe)=="0" goto wait
move /Y "B:\git\Boofutils\boofutils.exe_old" "%!s(MISSING)"
del "%!s(MISSING)"
