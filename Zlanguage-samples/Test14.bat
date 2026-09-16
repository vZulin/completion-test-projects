@rem
@rem Copyright 2015 the original author or authors.
@rem
@rem Licensed under the Apache License, Version 2.0 (the "License");
@rem you may not use this file except in compliance with the License.
@rem You may obtain a copy of the License at
@rem
@rem      https://www.apache.org/licenses/LICENSE-2.0
@rem
@rem Unless required by applicable law or agreed to in writing, software
@rem distributed under the License is distributed on an "AS IS" BASIS,
@rem WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@rem See the License for the specific language governing permissions and
@rem limitations under the License.
@rem
@rem SPDX-License-Identifier: Apache-2.0
@rem

@if "%DEBUG%"=="" @echo off
@rem ##########################################################################
@rem
@rem  tbcli startup script for Windows
@rem
@rem ##########################################################################

@rem Set local scope for the variables with windows NT shell
if "%OS%"=="Windows_NT" setlocal

set DIRNAME=%~dp0
if "%DIRNAME%"=="" set DIRNAME=.
@rem This is normally unused
set APP_BASE_NAME=%~n0
set APP_HOME=%DIRNAME%..

@rem Resolve any "." and ".." in APP_HOME to make it shorter.
for %%i in ("%APP_HOME%") do set APP_HOME=%%~fi

goto :skipFindBundledJava

:findBundledJava
if exist "%APP_HOME%\jre\bin\java.exe" (
    set "JAVA_EXE=%APP_HOME%\jre\bin\java.exe"
    exit /b 0
)
exit /b 1

:skipFindBundledJava

@rem Add default JVM options here. You can also use TB_JAVA_OPTS and TBCLI_OPTS to pass JVM options to this script.
set DEFAULT_JVM_OPTS="-ea" "-Dfile.encoding=UTF8" "-Xss384k" "-XX:+UnlockExperimentalVMOptions" "-XX:+CreateCoredumpOnCrash" "-XX:+UseStringDeduplication" "-XX:+UseCompressedOops" "-XX:+UseCompressedClassPointers" "-XX:MinMetaspaceFreeRatio=10" "-XX:MaxMetaspaceFreeRatio=10" "-XX:+UseSerialGC" "-XX:MinHeapFreeRatio=10" "-XX:MaxHeapFreeRatio=10" "-XX:-ShrinkHeapInSteps"

set JNA_LIBRARY_PATH="%APP_HOME%\lib;%APP_HOME%\lib\win32-aarch64;%APP_HOME%\lib\win32-x86-64"
set JNA_OPTS="-Djna.debug_load.jna=true" "-Djna.nosys=true" "-Djna.noclasspath=true" "-Djna.nounpack=true"

@rem Find java.exe
if defined TB_JAVA_HOME goto findJavaFromJavaHome

@rem Try bundled JRE first
call :findBundledJava
if %ERRORLEVEL% equ 0 goto execute

set JAVA_EXE=java.exe
%JAVA_EXE% -version >NUL 2>&1
if %ERRORLEVEL% equ 0 goto execute

echo. 1>&2
echo ERROR: TB_JAVA_HOME is not set and no 'java' command could be found in your PATH. 1>&2
echo No bundled JRE found in %APP_HOME%\jre 1>&2
echo. 1>&2
echo Please set the TB_JAVA_HOME variable in your environment to match the 1>&2
echo location of your Java installation. 1>&2

goto fail

:findJavaFromJavaHome
set TB_JAVA_HOME=%TB_JAVA_HOME:"=%
set JAVA_EXE=%TB_JAVA_HOME%/bin/java.exe

if exist "%JAVA_EXE%" goto execute

echo. 1>&2
echo ERROR: TB_JAVA_HOME is set to an invalid directory: %TB_JAVA_HOME% 1>&2
echo. 1>&2
echo Please set the TB_JAVA_HOME variable in your environment to match the 1>&2
echo location of your Java installation. 1>&2

goto fail

:execute
@rem Setup the command line

set CLASSPATH=%APP_HOME%/lib/*


set is_agent_proxy=false

:loop
if "%1"=="" goto end
if "%1"=="agent-proxy" (
    set is_agent_proxy=true
    goto end
)
shift
goto loop

:end

if "%is_agent_proxy%"=="true" (
    set DEFAULT_JVM_OPTS=%DEFAULT_JVM_OPTS% "-Xmx20m"
) else (
    set DEFAULT_JVM_OPTS=%DEFAULT_JVM_OPTS% "-Xmx120m"
)

if defined TB_DEBUG_PORT (
   echo java debug port: %TB_DEBUG_PORT%
   set debug_option=-Xdebug -Xrunjdwp:transport=dt_socket,server=y,suspend=%TB_DEBUG_SUSPEND:n%,address=%TB_DEBUG_PORT%
   set DEFAULT_JVM_OPTS=%DEFAULT_JVM_OPTS% "%debug_option%"
)

SETLOCAL ENABLEDELAYEDEXPANSION

if defined TB_ASYNC_PROFILER_PATH (
  if not defined TB_ASYNC_PROFILER_SNAPSHOTS_DIR (
    for %%f in ("%TB_CLI_PATH%") do set "TB_ASYNC_PROFILER_SNAPSHOTS_DIR=%%~dpf"
  )
  mkdir !TB_ASYNC_PROFILER_SNAPSHOTS_DIR!
  for /F %%i in ('powershell -NoProfile -Command "Get-Date -Format yyyyMMddHHmmss"') do set datetime=%%i
  set year=!datetime:~0,4!
  set month=!datetime:~4,2!
  set day=!datetime:~6,2!
  set hour=!datetime:~8,2!
  set minute=!datetime:~10,2!
  set second=!datetime:~12,2!
  set timestamp=!year!-!month!-!day!-!hour!-!minute!-!second!
  set snapshot_file=!TB_ASYNC_PROFILER_SNAPSHOTS_DIR!\!timestamp!.jfr
  set agent_param="-agentpath:%TB_ASYNC_PROFILER_PATH%=start,file=!snapshot_file!,event=wall,interval=10ms,jfr,jfrsync=profile"
  set DEFAULT_JVM_OPTS=!agent_param! %DEFAULT_JVM_OPTS%
)


@rem Execute tbcli
"%JAVA_EXE%" %DEFAULT_JVM_OPTS% %TB_JAVA_OPTS% %JNA_OPTS% %TBCLI_OPTS%  -Djna.boot.library.path="%JNA_LIBRARY_PATH%" -classpath "%CLASSPATH%" com.jetbrains.toolbox.MainKt %*

:end
@rem End local scope for the variables with windows NT shell
if %ERRORLEVEL% equ 0 goto mainEnd

:fail
rem Set variable TBCLI_EXIT_CONSOLE if you need the _script_ return code instead of
rem the _cmd.exe /c_ return code!
set EXIT_CODE=%ERRORLEVEL%
if %EXIT_CODE% equ 0 set EXIT_CODE=1
if not ""=="%TBCLI_EXIT_CONSOLE%" exit %EXIT_CODE%
exit /b %EXIT_CODE%

:mainEnd
if "%OS%"=="Windows_NT" endlocal

:omega
