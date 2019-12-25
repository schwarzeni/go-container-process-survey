#define _GNU_SOURCE /* See feature_test_macros(7) */
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>

__attribute__((constructor)) void enter_namespace(void)
{
  // const char *ACTIVE_FLAG_ENV = "mydocker_flag"; // 用于判断是否调用此函数
  const char *PID_ENV = "mydocker_pid";
  const char *CMD_ENV = "mydocker_cmd";
  // char *mydocker_flag = getenv(ACTIVE_FLAG_ENV);
  char *mydocker_pid = getenv(PID_ENV); // 进程需要进入的pid
  char *mydocker_cmd = getenv(CMD_ENV); // 需要执行的命令

  // if (!mydocker_flag) // 不需要执行此函数
  // {
  // exit(EXIT_SUCCESS);
  // }

  // check env
  if (!mydocker_pid)
  {
    //   fprintf(stderr, "missing %s env skip nsenter", PID_ENV);
    //   exit(EXIT_FAILURE);
    return;
  }
  if (!mydocker_cmd)
  {
    //   fprintf(stderr, "missing %s env skip nsenter", CMD_ENV);
    return;
    //   exit(EXIT_FAILURE);
  }

  char nspath[1024];
  char nslist[][5] = {"ipc", "uts", "net", "pid", "mnt"}; // 需要进入的五种namespace

  for (size_t i = 0; i < 5; i++)
  {
    sprintf(nspath, "/proc/%s/ns/%s", mydocker_pid, nslist[i]);
    int fd = open(nspath, O_RDONLY);
    // if (fd == -1)
    // {
    //   fprintf(stderr, "enter %s failed", nspath);
    //   exit(EXIT_FAILURE);
    // }
    if (setns(fd, 0) == -1)
    {
      fprintf(stderr, "enter ns %s failed", nspath);
      exit(EXIT_FAILURE);
    }
    if (close(fd) == -1)
    {
      fprintf(stderr, "close fd for %s failed", nspath);
      exit(EXIT_FAILURE);
    }
  }

  int res = system(mydocker_cmd);
  exit(EXIT_SUCCESS);
}
