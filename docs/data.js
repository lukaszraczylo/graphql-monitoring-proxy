window.BENCHMARK_DATA = {
  "lastUpdate": 1724063604894,
  "repoUrl": "https://github.com/lukaszraczylo/graphql-monitoring-proxy",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo",
            "email": "lukasz@raczylo.com"
          },
          "committer": {
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo",
            "email": "lukasz@raczylo.com"
          },
          "id": "4ca8ce57513ec2563b36eddaf83217082a108b98",
          "message": "fixup! fixup! fixup! Fix redis cache benchmark.",
          "timestamp": "2024-06-28T12:57:43Z",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/4ca8ce57513ec2563b36eddaf83217082a108b98"
        },
        "date": 1719579926085,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12909,
            "unit": "ns/op\t     228 B/op\t       5 allocs/op",
            "extra": "90420 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12909,
            "unit": "ns/op",
            "extra": "90420 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 228,
            "unit": "B/op",
            "extra": "90420 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "90420 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 676.4,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1737909 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 676.4,
            "unit": "ns/op",
            "extra": "1737909 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1737909 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1737909 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5248922,
            "unit": "ns/op\t  815883 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5248922,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815883,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000469,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000469,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13498,
            "unit": "ns/op\t     354 B/op\t      10 allocs/op",
            "extra": "88165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13498,
            "unit": "ns/op",
            "extra": "88165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "88165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "88165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 863,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1414516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 863,
            "unit": "ns/op",
            "extra": "1414516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1414516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1414516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13476,
            "unit": "ns/op\t     570 B/op\t      13 allocs/op",
            "extra": "91406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13476,
            "unit": "ns/op",
            "extra": "91406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 570,
            "unit": "B/op",
            "extra": "91406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "91406 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6364,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6364,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6278,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6278,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.626,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.626,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6186,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6186,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 838.5,
            "unit": "ns/op\t     344 B/op\t       4 allocs/op",
            "extra": "1480324 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 838.5,
            "unit": "ns/op",
            "extra": "1480324 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 344,
            "unit": "B/op",
            "extra": "1480324 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1480324 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 842.6,
            "unit": "ns/op\t     351 B/op\t       4 allocs/op",
            "extra": "1441490 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 842.6,
            "unit": "ns/op",
            "extra": "1441490 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "1441490 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1441490 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 840,
            "unit": "ns/op\t     348 B/op\t       4 allocs/op",
            "extra": "1457232 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 840,
            "unit": "ns/op",
            "extra": "1457232 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 348,
            "unit": "B/op",
            "extra": "1457232 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1457232 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 939.1,
            "unit": "ns/op\t     382 B/op\t       4 allocs/op",
            "extra": "1296472 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 939.1,
            "unit": "ns/op",
            "extra": "1296472 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 382,
            "unit": "B/op",
            "extra": "1296472 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1296472 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 856.5,
            "unit": "ns/op\t     351 B/op\t       4 allocs/op",
            "extra": "1440241 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 856.5,
            "unit": "ns/op",
            "extra": "1440241 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "1440241 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1440241 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1151,
            "unit": "ns/op\t     567 B/op\t       8 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1151,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 567,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 415.8,
            "unit": "ns/op\t     248 B/op\t       6 allocs/op",
            "extra": "3015525 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 415.8,
            "unit": "ns/op",
            "extra": "3015525 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "3015525 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "3015525 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 802.3,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1507239 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 802.3,
            "unit": "ns/op",
            "extra": "1507239 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1507239 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1507239 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b1ffffd545eacae0228f203f46736c85b8d616dc",
          "message": "Create static.yml",
          "timestamp": "2024-06-28T17:49:21+01:00",
          "tree_id": "23eaebc8cde5d657085b7fb615a0e3c22f8c0930",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/b1ffffd545eacae0228f203f46736c85b8d616dc"
        },
        "date": 1719593697908,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13019,
            "unit": "ns/op\t     219 B/op\t       5 allocs/op",
            "extra": "89265 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13019,
            "unit": "ns/op",
            "extra": "89265 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 219,
            "unit": "B/op",
            "extra": "89265 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "89265 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 694.6,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1720644 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 694.6,
            "unit": "ns/op",
            "extra": "1720644 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1720644 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1720644 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5265552,
            "unit": "ns/op\t  815877 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5265552,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815877,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000938,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000938,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13012,
            "unit": "ns/op\t     365 B/op\t      10 allocs/op",
            "extra": "86428 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13012,
            "unit": "ns/op",
            "extra": "86428 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 365,
            "unit": "B/op",
            "extra": "86428 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "86428 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 860.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1372743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 860.6,
            "unit": "ns/op",
            "extra": "1372743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1372743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1372743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13225,
            "unit": "ns/op\t     560 B/op\t      13 allocs/op",
            "extra": "95506 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13225,
            "unit": "ns/op",
            "extra": "95506 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "95506 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "95506 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.621,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.621,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6245,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6245,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6237,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6237,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6188,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6188,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 833.7,
            "unit": "ns/op\t     345 B/op\t       4 allocs/op",
            "extra": "1471714 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 833.7,
            "unit": "ns/op",
            "extra": "1471714 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 345,
            "unit": "B/op",
            "extra": "1471714 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1471714 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 843.2,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1444736 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 843.2,
            "unit": "ns/op",
            "extra": "1444736 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1444736 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1444736 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 848.1,
            "unit": "ns/op\t     351 B/op\t       4 allocs/op",
            "extra": "1441908 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 848.1,
            "unit": "ns/op",
            "extra": "1441908 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "1441908 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1441908 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 968.2,
            "unit": "ns/op\t     381 B/op\t       4 allocs/op",
            "extra": "1300447 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 968.2,
            "unit": "ns/op",
            "extra": "1300447 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 381,
            "unit": "B/op",
            "extra": "1300447 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1300447 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 882.2,
            "unit": "ns/op\t     358 B/op\t       4 allocs/op",
            "extra": "1406826 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 882.2,
            "unit": "ns/op",
            "extra": "1406826 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 358,
            "unit": "B/op",
            "extra": "1406826 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1406826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1156,
            "unit": "ns/op\t     567 B/op\t       8 allocs/op",
            "extra": "992169 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1156,
            "unit": "ns/op",
            "extra": "992169 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 567,
            "unit": "B/op",
            "extra": "992169 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "992169 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 399,
            "unit": "ns/op\t     248 B/op\t       6 allocs/op",
            "extra": "2981295 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 399,
            "unit": "ns/op",
            "extra": "2981295 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "2981295 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2981295 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 796.6,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1552028 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 796.6,
            "unit": "ns/op",
            "extra": "1552028 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1552028 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1552028 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "distinct": true,
          "id": "162c4acd7cd63d4cf96ecb3142467cb4fff8198f",
          "message": "fixup! fixup! fixup! fixup! fixup! fixup! Fix redis cache benchmark.",
          "timestamp": "2024-06-28T18:05:17+01:00",
          "tree_id": "f9717f06910488d6035ed840f8fc807b3f55e251",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/162c4acd7cd63d4cf96ecb3142467cb4fff8198f"
        },
        "date": 1719594657369,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12516,
            "unit": "ns/op\t     219 B/op\t       5 allocs/op",
            "extra": "89520 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12516,
            "unit": "ns/op",
            "extra": "89520 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 219,
            "unit": "B/op",
            "extra": "89520 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "89520 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 686.5,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1756580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 686.5,
            "unit": "ns/op",
            "extra": "1756580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1756580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1756580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5239163,
            "unit": "ns/op\t  815895 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5239163,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815895,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000504,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000504,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12500,
            "unit": "ns/op\t     349 B/op\t      10 allocs/op",
            "extra": "93064 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12500,
            "unit": "ns/op",
            "extra": "93064 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 349,
            "unit": "B/op",
            "extra": "93064 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "93064 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 851.8,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1455595 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 851.8,
            "unit": "ns/op",
            "extra": "1455595 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1455595 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1455595 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13590,
            "unit": "ns/op\t     568 B/op\t      13 allocs/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13590,
            "unit": "ns/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 568,
            "unit": "B/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6207,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6249,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6249,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6202,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6202,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6223,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6223,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 838.3,
            "unit": "ns/op\t     346 B/op\t       4 allocs/op",
            "extra": "1468740 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 838.3,
            "unit": "ns/op",
            "extra": "1468740 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 346,
            "unit": "B/op",
            "extra": "1468740 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1468740 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 857,
            "unit": "ns/op\t     347 B/op\t       4 allocs/op",
            "extra": "1463252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 857,
            "unit": "ns/op",
            "extra": "1463252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 347,
            "unit": "B/op",
            "extra": "1463252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1463252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 850.9,
            "unit": "ns/op\t     351 B/op\t       4 allocs/op",
            "extra": "1438802 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 850.9,
            "unit": "ns/op",
            "extra": "1438802 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "1438802 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1438802 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 959.1,
            "unit": "ns/op\t     381 B/op\t       4 allocs/op",
            "extra": "1299824 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 959.1,
            "unit": "ns/op",
            "extra": "1299824 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 381,
            "unit": "B/op",
            "extra": "1299824 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1299824 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 857.1,
            "unit": "ns/op\t     356 B/op\t       4 allocs/op",
            "extra": "1415732 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 857.1,
            "unit": "ns/op",
            "extra": "1415732 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 356,
            "unit": "B/op",
            "extra": "1415732 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1415732 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1138,
            "unit": "ns/op\t     567 B/op\t       8 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1138,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 567,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 399.3,
            "unit": "ns/op\t     248 B/op\t       6 allocs/op",
            "extra": "3002569 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 399.3,
            "unit": "ns/op",
            "extra": "3002569 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 248,
            "unit": "B/op",
            "extra": "3002569 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "3002569 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 773.4,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1512163 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 773.4,
            "unit": "ns/op",
            "extra": "1512163 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1512163 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1512163 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d141fe3c041a7e6d2dc1e08feb53669908f82ef0",
          "message": "Fix the introduced bug where RO endpoint could've been accidentally used. (#17)\n\n* Fix the introduced bug where RO endpoint could've been accidentally used.",
          "timestamp": "2024-06-28T21:48:39+01:00",
          "tree_id": "b5e14a3b28655edc70b4d9f1966de28446931622",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/d141fe3c041a7e6d2dc1e08feb53669908f82ef0"
        },
        "date": 1719608039170,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 434.6,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2312928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 434.6,
            "unit": "ns/op",
            "extra": "2312928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2312928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2312928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8699,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "126346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8699,
            "unit": "ns/op",
            "extra": "126346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "126346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "126346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5918,
            "unit": "ns/op\t     655 B/op\t       8 allocs/op",
            "extra": "181908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5918,
            "unit": "ns/op",
            "extra": "181908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 655,
            "unit": "B/op",
            "extra": "181908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "181908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9662,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "121544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9662,
            "unit": "ns/op",
            "extra": "121544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "121544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "121544 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13299,
            "unit": "ns/op\t     210 B/op\t       5 allocs/op",
            "extra": "88940 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13299,
            "unit": "ns/op",
            "extra": "88940 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 210,
            "unit": "B/op",
            "extra": "88940 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "88940 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 674.3,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1719891 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 674.3,
            "unit": "ns/op",
            "extra": "1719891 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1719891 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1719891 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5250908,
            "unit": "ns/op\t  815893 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5250908,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815893,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000717,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000717,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 14016,
            "unit": "ns/op\t     352 B/op\t      10 allocs/op",
            "extra": "90498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 14016,
            "unit": "ns/op",
            "extra": "90498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "90498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "90498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 868.7,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1413685 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 868.7,
            "unit": "ns/op",
            "extra": "1413685 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1413685 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1413685 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 14094,
            "unit": "ns/op\t     564 B/op\t      13 allocs/op",
            "extra": "83424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 14094,
            "unit": "ns/op",
            "extra": "83424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 564,
            "unit": "B/op",
            "extra": "83424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "83424 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6285,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6285,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6326,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6326,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6226,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6226,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6192,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6192,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 849.4,
            "unit": "ns/op\t     346 B/op\t       4 allocs/op",
            "extra": "1465444 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 849.4,
            "unit": "ns/op",
            "extra": "1465444 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 346,
            "unit": "B/op",
            "extra": "1465444 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1465444 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 842.7,
            "unit": "ns/op\t     349 B/op\t       4 allocs/op",
            "extra": "1450833 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 842.7,
            "unit": "ns/op",
            "extra": "1450833 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 349,
            "unit": "B/op",
            "extra": "1450833 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1450833 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 859,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1411665 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 859,
            "unit": "ns/op",
            "extra": "1411665 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1411665 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1411665 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 953.4,
            "unit": "ns/op\t     384 B/op\t       4 allocs/op",
            "extra": "1286743 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 953.4,
            "unit": "ns/op",
            "extra": "1286743 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "1286743 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1286743 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 863.1,
            "unit": "ns/op\t     356 B/op\t       4 allocs/op",
            "extra": "1415925 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 863.1,
            "unit": "ns/op",
            "extra": "1415925 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 356,
            "unit": "B/op",
            "extra": "1415925 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1415925 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1314,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "925789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1314,
            "unit": "ns/op",
            "extra": "925789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "925789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "925789 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 493.8,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2395348 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 493.8,
            "unit": "ns/op",
            "extra": "2395348 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2395348 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2395348 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 786.9,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1528592 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 786.9,
            "unit": "ns/op",
            "extra": "1528592 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1528592 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1528592 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "distinct": true,
          "id": "a24e6c8c4dbacf70636a6fc3622f27d620c3a7b3",
          "message": "fixup! Fix the introduced bug where RO endpoint could've been accidentally used. (#17)",
          "timestamp": "2024-06-29T08:52:41+01:00",
          "tree_id": "6e7538aab35d433e99825fd68727d7420780c584",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/a24e6c8c4dbacf70636a6fc3622f27d620c3a7b3"
        },
        "date": 1719647892840,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 434,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2543504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 434,
            "unit": "ns/op",
            "extra": "2543504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2543504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2543504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8596,
            "unit": "ns/op\t     501 B/op\t      23 allocs/op",
            "extra": "138549 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8596,
            "unit": "ns/op",
            "extra": "138549 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 501,
            "unit": "B/op",
            "extra": "138549 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "138549 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5824,
            "unit": "ns/op\t     695 B/op\t       8 allocs/op",
            "extra": "182721 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5824,
            "unit": "ns/op",
            "extra": "182721 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 695,
            "unit": "B/op",
            "extra": "182721 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "182721 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9652,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "121965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9652,
            "unit": "ns/op",
            "extra": "121965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "121965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "121965 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12806,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "98402 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12806,
            "unit": "ns/op",
            "extra": "98402 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "98402 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "98402 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 679.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1745907 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 679.6,
            "unit": "ns/op",
            "extra": "1745907 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1745907 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1745907 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5249062,
            "unit": "ns/op\t  815891 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5249062,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815891,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000467,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000467,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13203,
            "unit": "ns/op\t     368 B/op\t      10 allocs/op",
            "extra": "91802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13203,
            "unit": "ns/op",
            "extra": "91802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 368,
            "unit": "B/op",
            "extra": "91802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "91802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 857,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1422978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 857,
            "unit": "ns/op",
            "extra": "1422978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1422978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1422978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13374,
            "unit": "ns/op\t     559 B/op\t      13 allocs/op",
            "extra": "100851 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13374,
            "unit": "ns/op",
            "extra": "100851 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 559,
            "unit": "B/op",
            "extra": "100851 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "100851 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6213,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6213,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6204,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6204,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6194,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6194,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6196,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6196,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 827.9,
            "unit": "ns/op\t     344 B/op\t       4 allocs/op",
            "extra": "1478889 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 827.9,
            "unit": "ns/op",
            "extra": "1478889 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 344,
            "unit": "B/op",
            "extra": "1478889 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1478889 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 844.5,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1446111 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 844.5,
            "unit": "ns/op",
            "extra": "1446111 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1446111 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1446111 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 851.6,
            "unit": "ns/op\t     353 B/op\t       4 allocs/op",
            "extra": "1431538 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 851.6,
            "unit": "ns/op",
            "extra": "1431538 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 353,
            "unit": "B/op",
            "extra": "1431538 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1431538 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 959.2,
            "unit": "ns/op\t     395 B/op\t       4 allocs/op",
            "extra": "1244581 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 959.2,
            "unit": "ns/op",
            "extra": "1244581 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 395,
            "unit": "B/op",
            "extra": "1244581 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1244581 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 873.8,
            "unit": "ns/op\t     359 B/op\t       4 allocs/op",
            "extra": "1401946 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 873.8,
            "unit": "ns/op",
            "extra": "1401946 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 359,
            "unit": "B/op",
            "extra": "1401946 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1401946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1296,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "936098 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1296,
            "unit": "ns/op",
            "extra": "936098 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "936098 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "936098 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 491.2,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2430316 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 491.2,
            "unit": "ns/op",
            "extra": "2430316 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2430316 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2430316 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 769.9,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1572525 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 769.9,
            "unit": "ns/op",
            "extra": "1572525 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1572525 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1572525 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "distinct": true,
          "id": "dfd3b02014443a1aae9af7d42c2544567eb6e8ec",
          "message": "Release 0.19.x",
          "timestamp": "2024-06-29T09:57:52+01:00",
          "tree_id": "9e0c4115395f8d9cb65450cb5989acd9e0ab9cb0",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/dfd3b02014443a1aae9af7d42c2544567eb6e8ec"
        },
        "date": 1719651833177,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 434,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2778037 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 434,
            "unit": "ns/op",
            "extra": "2778037 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2778037 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2778037 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8697,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "126434 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8697,
            "unit": "ns/op",
            "extra": "126434 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "126434 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "126434 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 6560,
            "unit": "ns/op\t     634 B/op\t       8 allocs/op",
            "extra": "179163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 6560,
            "unit": "ns/op",
            "extra": "179163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 634,
            "unit": "B/op",
            "extra": "179163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "179163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9676,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "123898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9676,
            "unit": "ns/op",
            "extra": "123898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "123898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "123898 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12535,
            "unit": "ns/op\t     210 B/op\t       5 allocs/op",
            "extra": "89528 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12535,
            "unit": "ns/op",
            "extra": "89528 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 210,
            "unit": "B/op",
            "extra": "89528 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "89528 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 682.4,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1714333 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 682.4,
            "unit": "ns/op",
            "extra": "1714333 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1714333 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1714333 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5245807,
            "unit": "ns/op\t  815883 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5245807,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815883,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000844,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000844,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13316,
            "unit": "ns/op\t     355 B/op\t      10 allocs/op",
            "extra": "87300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13316,
            "unit": "ns/op",
            "extra": "87300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "87300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "87300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 914.9,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1406504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 914.9,
            "unit": "ns/op",
            "extra": "1406504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1406504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1406504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12864,
            "unit": "ns/op\t     561 B/op\t      13 allocs/op",
            "extra": "91402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12864,
            "unit": "ns/op",
            "extra": "91402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "91402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "91402 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6289,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6289,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6398,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6398,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6245,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6245,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6239,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6239,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 839.2,
            "unit": "ns/op\t     345 B/op\t       4 allocs/op",
            "extra": "1471993 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 839.2,
            "unit": "ns/op",
            "extra": "1471993 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 345,
            "unit": "B/op",
            "extra": "1471993 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1471993 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 853.9,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1447663 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 853.9,
            "unit": "ns/op",
            "extra": "1447663 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1447663 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1447663 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 859.3,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1412090 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 859.3,
            "unit": "ns/op",
            "extra": "1412090 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1412090 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1412090 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 962.3,
            "unit": "ns/op\t     385 B/op\t       4 allocs/op",
            "extra": "1286086 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 962.3,
            "unit": "ns/op",
            "extra": "1286086 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 385,
            "unit": "B/op",
            "extra": "1286086 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1286086 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 873,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1408712 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 873,
            "unit": "ns/op",
            "extra": "1408712 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1408712 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1408712 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1278,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "900110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1278,
            "unit": "ns/op",
            "extra": "900110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "900110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "900110 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 491.6,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2431520 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 491.6,
            "unit": "ns/op",
            "extra": "2431520 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2431520 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2431520 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 793.6,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1496040 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 793.6,
            "unit": "ns/op",
            "extra": "1496040 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1496040 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1496040 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "8bdc151c7e58f75ddbffd656b1b32b3202e057cd",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-06-30T03:02:53Z",
          "tree_id": "295464bed478a5fa2c2b7b2393ac79bf8e5512e6",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/8bdc151c7e58f75ddbffd656b1b32b3202e057cd"
        },
        "date": 1719716982103,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 431.9,
            "unit": "ns/op\t     562 B/op\t       2 allocs/op",
            "extra": "2778376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 431.9,
            "unit": "ns/op",
            "extra": "2778376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 562,
            "unit": "B/op",
            "extra": "2778376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2778376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8670,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "118795 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8670,
            "unit": "ns/op",
            "extra": "118795 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "118795 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "118795 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5836,
            "unit": "ns/op\t     635 B/op\t       8 allocs/op",
            "extra": "179043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5836,
            "unit": "ns/op",
            "extra": "179043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 635,
            "unit": "B/op",
            "extra": "179043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "179043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9595,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "122193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9595,
            "unit": "ns/op",
            "extra": "122193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "122193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "122193 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12136,
            "unit": "ns/op\t     217 B/op\t       5 allocs/op",
            "extra": "96874 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12136,
            "unit": "ns/op",
            "extra": "96874 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 217,
            "unit": "B/op",
            "extra": "96874 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "96874 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 694.4,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1679307 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 694.4,
            "unit": "ns/op",
            "extra": "1679307 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1679307 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1679307 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5255388,
            "unit": "ns/op\t  815882 B/op\t      36 allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5255388,
            "unit": "ns/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815882,
            "unit": "B/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000861,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000861,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13155,
            "unit": "ns/op\t     348 B/op\t      10 allocs/op",
            "extra": "95116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13155,
            "unit": "ns/op",
            "extra": "95116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 348,
            "unit": "B/op",
            "extra": "95116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "95116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 852.5,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1358151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 852.5,
            "unit": "ns/op",
            "extra": "1358151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1358151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1358151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13370,
            "unit": "ns/op\t     553 B/op\t      13 allocs/op",
            "extra": "90700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13370,
            "unit": "ns/op",
            "extra": "90700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 553,
            "unit": "B/op",
            "extra": "90700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "90700 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6331,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6331,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6363,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6363,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6196,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6196,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6256,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6256,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 835.4,
            "unit": "ns/op\t     343 B/op\t       4 allocs/op",
            "extra": "1480497 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 835.4,
            "unit": "ns/op",
            "extra": "1480497 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 343,
            "unit": "B/op",
            "extra": "1480497 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1480497 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 835.6,
            "unit": "ns/op\t     348 B/op\t       4 allocs/op",
            "extra": "1454702 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 835.6,
            "unit": "ns/op",
            "extra": "1454702 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 348,
            "unit": "B/op",
            "extra": "1454702 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1454702 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 852.2,
            "unit": "ns/op\t     353 B/op\t       4 allocs/op",
            "extra": "1431334 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 852.2,
            "unit": "ns/op",
            "extra": "1431334 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 353,
            "unit": "B/op",
            "extra": "1431334 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1431334 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 953.7,
            "unit": "ns/op\t     383 B/op\t       4 allocs/op",
            "extra": "1292127 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 953.7,
            "unit": "ns/op",
            "extra": "1292127 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 383,
            "unit": "B/op",
            "extra": "1292127 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1292127 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 864.8,
            "unit": "ns/op\t     355 B/op\t       4 allocs/op",
            "extra": "1421230 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 864.8,
            "unit": "ns/op",
            "extra": "1421230 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "1421230 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1421230 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1306,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "924477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1306,
            "unit": "ns/op",
            "extra": "924477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "924477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "924477 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 502.2,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2285773 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 502.2,
            "unit": "ns/op",
            "extra": "2285773 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2285773 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2285773 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 808.8,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1488075 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 808.8,
            "unit": "ns/op",
            "extra": "1488075 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1488075 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1488075 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "e28da35ca4053e64c5e78027f855b4ebc29dfc9e",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-01T03:01:35Z",
          "tree_id": "3c60382e1795d4d8a17ac2e41fd1bce1dc1ef9fe",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/e28da35ca4053e64c5e78027f855b4ebc29dfc9e"
        },
        "date": 1719803295226,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 430.1,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2761239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 430.1,
            "unit": "ns/op",
            "extra": "2761239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2761239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2761239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8496,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "125218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8496,
            "unit": "ns/op",
            "extra": "125218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "125218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "125218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5865,
            "unit": "ns/op\t     648 B/op\t       8 allocs/op",
            "extra": "194928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5865,
            "unit": "ns/op",
            "extra": "194928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 648,
            "unit": "B/op",
            "extra": "194928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "194928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9620,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "121251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9620,
            "unit": "ns/op",
            "extra": "121251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "121251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "121251 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12586,
            "unit": "ns/op\t     201 B/op\t       5 allocs/op",
            "extra": "90541 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12586,
            "unit": "ns/op",
            "extra": "90541 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 201,
            "unit": "B/op",
            "extra": "90541 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "90541 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 687.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1789998 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 687.6,
            "unit": "ns/op",
            "extra": "1789998 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1789998 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1789998 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5249393,
            "unit": "ns/op\t  815883 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5249393,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815883,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000529,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000529,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13213,
            "unit": "ns/op\t     362 B/op\t      10 allocs/op",
            "extra": "89137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13213,
            "unit": "ns/op",
            "extra": "89137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 362,
            "unit": "B/op",
            "extra": "89137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "89137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 853.5,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1397914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 853.5,
            "unit": "ns/op",
            "extra": "1397914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1397914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1397914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13436,
            "unit": "ns/op\t     551 B/op\t      13 allocs/op",
            "extra": "98412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13436,
            "unit": "ns/op",
            "extra": "98412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 551,
            "unit": "B/op",
            "extra": "98412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "98412 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6044,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6044,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6173,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6173,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6125,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6125,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6187,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6187,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 814.1,
            "unit": "ns/op\t     345 B/op\t       4 allocs/op",
            "extra": "1473991 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 814.1,
            "unit": "ns/op",
            "extra": "1473991 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 345,
            "unit": "B/op",
            "extra": "1473991 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1473991 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 799.8,
            "unit": "ns/op\t     345 B/op\t       4 allocs/op",
            "extra": "1470840 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 799.8,
            "unit": "ns/op",
            "extra": "1470840 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 345,
            "unit": "B/op",
            "extra": "1470840 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1470840 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 814.4,
            "unit": "ns/op\t     342 B/op\t       4 allocs/op",
            "extra": "1490410 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 814.4,
            "unit": "ns/op",
            "extra": "1490410 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 342,
            "unit": "B/op",
            "extra": "1490410 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1490410 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 956.9,
            "unit": "ns/op\t     383 B/op\t       4 allocs/op",
            "extra": "1293582 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 956.9,
            "unit": "ns/op",
            "extra": "1293582 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 383,
            "unit": "B/op",
            "extra": "1293582 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1293582 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 842.5,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1448175 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 842.5,
            "unit": "ns/op",
            "extra": "1448175 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1448175 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1448175 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1304,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "935335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1304,
            "unit": "ns/op",
            "extra": "935335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "935335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "935335 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 483.4,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2448729 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 483.4,
            "unit": "ns/op",
            "extra": "2448729 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2448729 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2448729 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 763.7,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1578711 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 763.7,
            "unit": "ns/op",
            "extra": "1578711 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1578711 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1578711 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "cb862ae4b1e2d0e5d67f9e8f7ed1bbce9ec13635",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-03T03:01:32Z",
          "tree_id": "6245faa3252919ca8705d6a9644907864c7f7a8c",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/cb862ae4b1e2d0e5d67f9e8f7ed1bbce9ec13635"
        },
        "date": 1719976089214,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 432.3,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2752399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 432.3,
            "unit": "ns/op",
            "extra": "2752399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2752399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2752399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8813,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "122277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8813,
            "unit": "ns/op",
            "extra": "122277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "122277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "122277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5756,
            "unit": "ns/op\t     662 B/op\t       8 allocs/op",
            "extra": "185430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5756,
            "unit": "ns/op",
            "extra": "185430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 662,
            "unit": "B/op",
            "extra": "185430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "185430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9640,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "123042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9640,
            "unit": "ns/op",
            "extra": "123042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "123042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "123042 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 11859,
            "unit": "ns/op\t     200 B/op\t       5 allocs/op",
            "extra": "97363 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 11859,
            "unit": "ns/op",
            "extra": "97363 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 200,
            "unit": "B/op",
            "extra": "97363 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "97363 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 677.1,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1759860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 677.1,
            "unit": "ns/op",
            "extra": "1759860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1759860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1759860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5253224,
            "unit": "ns/op\t  815889 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5253224,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815889,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000462,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000462,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12281,
            "unit": "ns/op\t     354 B/op\t      10 allocs/op",
            "extra": "97904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12281,
            "unit": "ns/op",
            "extra": "97904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "97904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "97904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 844.2,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1405993 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 844.2,
            "unit": "ns/op",
            "extra": "1405993 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1405993 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1405993 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13161,
            "unit": "ns/op\t     552 B/op\t      13 allocs/op",
            "extra": "93609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13161,
            "unit": "ns/op",
            "extra": "93609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "93609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "93609 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.63,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6421,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6421,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6191,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6191,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6218,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6218,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 837.5,
            "unit": "ns/op\t     344 B/op\t       4 allocs/op",
            "extra": "1477418 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 837.5,
            "unit": "ns/op",
            "extra": "1477418 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 344,
            "unit": "B/op",
            "extra": "1477418 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1477418 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 842.3,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1443800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 842.3,
            "unit": "ns/op",
            "extra": "1443800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1443800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1443800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 851.6,
            "unit": "ns/op\t     353 B/op\t       4 allocs/op",
            "extra": "1432281 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 851.6,
            "unit": "ns/op",
            "extra": "1432281 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 353,
            "unit": "B/op",
            "extra": "1432281 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1432281 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 948.8,
            "unit": "ns/op\t     381 B/op\t       4 allocs/op",
            "extra": "1299738 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 948.8,
            "unit": "ns/op",
            "extra": "1299738 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 381,
            "unit": "B/op",
            "extra": "1299738 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1299738 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 864.7,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1396107 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 864.7,
            "unit": "ns/op",
            "extra": "1396107 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1396107 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1396107 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1296,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "938977 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1296,
            "unit": "ns/op",
            "extra": "938977 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "938977 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "938977 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 491.2,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2430802 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 491.2,
            "unit": "ns/op",
            "extra": "2430802 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2430802 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2430802 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 798.5,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1531870 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 798.5,
            "unit": "ns/op",
            "extra": "1531870 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1531870 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1531870 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "a2986dfc1ad5f792a3d79e14712f382ded24875c",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-05T03:02:00Z",
          "tree_id": "f8e77b82b3129d53f1b90639cdc0f2511a01d0ec",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/a2986dfc1ad5f792a3d79e14712f382ded24875c"
        },
        "date": 1720148952545,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 438.4,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2309494 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 438.4,
            "unit": "ns/op",
            "extra": "2309494 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2309494 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2309494 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8756,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "126268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8756,
            "unit": "ns/op",
            "extra": "126268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "126268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "126268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5914,
            "unit": "ns/op\t     624 B/op\t       8 allocs/op",
            "extra": "182421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5914,
            "unit": "ns/op",
            "extra": "182421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 624,
            "unit": "B/op",
            "extra": "182421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "182421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9838,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "118530 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9838,
            "unit": "ns/op",
            "extra": "118530 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "118530 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "118530 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12283,
            "unit": "ns/op\t     217 B/op\t       5 allocs/op",
            "extra": "94934 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12283,
            "unit": "ns/op",
            "extra": "94934 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 217,
            "unit": "B/op",
            "extra": "94934 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "94934 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 700,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1685227 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 700,
            "unit": "ns/op",
            "extra": "1685227 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1685227 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1685227 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5241500,
            "unit": "ns/op\t  815888 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5241500,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815888,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000965,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000965,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12657,
            "unit": "ns/op\t     349 B/op\t      10 allocs/op",
            "extra": "93811 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12657,
            "unit": "ns/op",
            "extra": "93811 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 349,
            "unit": "B/op",
            "extra": "93811 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "93811 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 888.3,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1416740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 888.3,
            "unit": "ns/op",
            "extra": "1416740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1416740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1416740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13914,
            "unit": "ns/op\t     568 B/op\t      13 allocs/op",
            "extra": "98557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13914,
            "unit": "ns/op",
            "extra": "98557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 568,
            "unit": "B/op",
            "extra": "98557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "98557 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6239,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6239,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6348,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6348,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6205,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6205,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6257,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6257,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 866.7,
            "unit": "ns/op\t     347 B/op\t       4 allocs/op",
            "extra": "1462238 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 866.7,
            "unit": "ns/op",
            "extra": "1462238 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 347,
            "unit": "B/op",
            "extra": "1462238 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1462238 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 848.2,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1435327 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 848.2,
            "unit": "ns/op",
            "extra": "1435327 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1435327 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1435327 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 839.8,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1409810 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 839.8,
            "unit": "ns/op",
            "extra": "1409810 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1409810 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1409810 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 962.6,
            "unit": "ns/op\t     383 B/op\t       4 allocs/op",
            "extra": "1293567 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 962.6,
            "unit": "ns/op",
            "extra": "1293567 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 383,
            "unit": "B/op",
            "extra": "1293567 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1293567 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 863,
            "unit": "ns/op\t     355 B/op\t       4 allocs/op",
            "extra": "1419394 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 863,
            "unit": "ns/op",
            "extra": "1419394 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "1419394 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1419394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1333,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "912627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1333,
            "unit": "ns/op",
            "extra": "912627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "912627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "912627 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 494.2,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2440372 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 494.2,
            "unit": "ns/op",
            "extra": "2440372 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2440372 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2440372 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 794.9,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1526533 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 794.9,
            "unit": "ns/op",
            "extra": "1526533 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1526533 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1526533 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "ab703d331ef791e9d1bd5393c974370506c8826b",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-06T03:01:42Z",
          "tree_id": "7c3cfc3e484ef0e28a0ac5d1e54e6719085f40c4",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/ab703d331ef791e9d1bd5393c974370506c8826b"
        },
        "date": 1720235300729,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 468.7,
            "unit": "ns/op\t     562 B/op\t       2 allocs/op",
            "extra": "2748195 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 468.7,
            "unit": "ns/op",
            "extra": "2748195 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 562,
            "unit": "B/op",
            "extra": "2748195 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2748195 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8853,
            "unit": "ns/op\t     501 B/op\t      23 allocs/op",
            "extra": "132964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8853,
            "unit": "ns/op",
            "extra": "132964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 501,
            "unit": "B/op",
            "extra": "132964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "132964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5942,
            "unit": "ns/op\t     690 B/op\t       8 allocs/op",
            "extra": "178087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5942,
            "unit": "ns/op",
            "extra": "178087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 690,
            "unit": "B/op",
            "extra": "178087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "178087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9711,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "120338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9711,
            "unit": "ns/op",
            "extra": "120338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "120338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "120338 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13446,
            "unit": "ns/op\t     210 B/op\t       5 allocs/op",
            "extra": "87766 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13446,
            "unit": "ns/op",
            "extra": "87766 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 210,
            "unit": "B/op",
            "extra": "87766 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "87766 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 673.5,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1724179 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 673.5,
            "unit": "ns/op",
            "extra": "1724179 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1724179 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1724179 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5224982,
            "unit": "ns/op\t  815870 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5224982,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815870,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000511,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000511,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13752,
            "unit": "ns/op\t     365 B/op\t      10 allocs/op",
            "extra": "86311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13752,
            "unit": "ns/op",
            "extra": "86311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 365,
            "unit": "B/op",
            "extra": "86311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "86311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 831.1,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1436809 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 831.1,
            "unit": "ns/op",
            "extra": "1436809 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1436809 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1436809 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 14759,
            "unit": "ns/op\t     553 B/op\t      13 allocs/op",
            "extra": "86416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 14759,
            "unit": "ns/op",
            "extra": "86416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 553,
            "unit": "B/op",
            "extra": "86416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "86416 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6234,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6234,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6351,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6351,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.623,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.623,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6191,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6191,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 832.7,
            "unit": "ns/op\t     345 B/op\t       4 allocs/op",
            "extra": "1473988 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 832.7,
            "unit": "ns/op",
            "extra": "1473988 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 345,
            "unit": "B/op",
            "extra": "1473988 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1473988 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 853,
            "unit": "ns/op\t     349 B/op\t       4 allocs/op",
            "extra": "1452800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 853,
            "unit": "ns/op",
            "extra": "1452800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 349,
            "unit": "B/op",
            "extra": "1452800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1452800 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 840.7,
            "unit": "ns/op\t     351 B/op\t       4 allocs/op",
            "extra": "1440433 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 840.7,
            "unit": "ns/op",
            "extra": "1440433 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "1440433 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1440433 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 936.3,
            "unit": "ns/op\t     378 B/op\t       4 allocs/op",
            "extra": "1314502 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 936.3,
            "unit": "ns/op",
            "extra": "1314502 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 378,
            "unit": "B/op",
            "extra": "1314502 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1314502 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 862.1,
            "unit": "ns/op\t     359 B/op\t       4 allocs/op",
            "extra": "1400917 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 862.1,
            "unit": "ns/op",
            "extra": "1400917 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 359,
            "unit": "B/op",
            "extra": "1400917 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1400917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1302,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "922354 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1302,
            "unit": "ns/op",
            "extra": "922354 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "922354 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "922354 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 493.3,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2445372 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 493.3,
            "unit": "ns/op",
            "extra": "2445372 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2445372 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2445372 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 791.5,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1478878 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 791.5,
            "unit": "ns/op",
            "extra": "1478878 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1478878 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1478878 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "306139fcef9c79cfb952f4cbb28627423a089382",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-10T03:01:46Z",
          "tree_id": "6ae3f425134059896f8b2a8a246570d6989481cd",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/306139fcef9c79cfb952f4cbb28627423a089382"
        },
        "date": 1720580903340,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 437.4,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2688733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 437.4,
            "unit": "ns/op",
            "extra": "2688733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2688733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2688733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8597,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "122214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8597,
            "unit": "ns/op",
            "extra": "122214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "122214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "122214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5875,
            "unit": "ns/op\t     695 B/op\t       8 allocs/op",
            "extra": "183080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5875,
            "unit": "ns/op",
            "extra": "183080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 695,
            "unit": "B/op",
            "extra": "183080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "183080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9459,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "120656 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9459,
            "unit": "ns/op",
            "extra": "120656 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "120656 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "120656 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12307,
            "unit": "ns/op\t     210 B/op\t       5 allocs/op",
            "extra": "89142 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12307,
            "unit": "ns/op",
            "extra": "89142 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 210,
            "unit": "B/op",
            "extra": "89142 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "89142 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 690,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1743315 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 690,
            "unit": "ns/op",
            "extra": "1743315 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1743315 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1743315 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5250943,
            "unit": "ns/op\t  815881 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5250943,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815881,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000954,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000954,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13748,
            "unit": "ns/op\t     360 B/op\t      10 allocs/op",
            "extra": "91464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13748,
            "unit": "ns/op",
            "extra": "91464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "91464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "91464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 861.3,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1394798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 861.3,
            "unit": "ns/op",
            "extra": "1394798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1394798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1394798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13239,
            "unit": "ns/op\t     562 B/op\t      13 allocs/op",
            "extra": "88581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13239,
            "unit": "ns/op",
            "extra": "88581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 562,
            "unit": "B/op",
            "extra": "88581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "88581 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.629,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.629,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6238,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6238,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6226,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6226,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6213,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6213,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 842.7,
            "unit": "ns/op\t     344 B/op\t       4 allocs/op",
            "extra": "1476423 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 842.7,
            "unit": "ns/op",
            "extra": "1476423 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 344,
            "unit": "B/op",
            "extra": "1476423 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1476423 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 847.5,
            "unit": "ns/op\t     351 B/op\t       4 allocs/op",
            "extra": "1441766 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 847.5,
            "unit": "ns/op",
            "extra": "1441766 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "1441766 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1441766 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 853.7,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1436352 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 853.7,
            "unit": "ns/op",
            "extra": "1436352 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1436352 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1436352 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 978.8,
            "unit": "ns/op\t     383 B/op\t       4 allocs/op",
            "extra": "1291178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 978.8,
            "unit": "ns/op",
            "extra": "1291178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 383,
            "unit": "B/op",
            "extra": "1291178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1291178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 882.2,
            "unit": "ns/op\t     356 B/op\t       4 allocs/op",
            "extra": "1413902 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 882.2,
            "unit": "ns/op",
            "extra": "1413902 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 356,
            "unit": "B/op",
            "extra": "1413902 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1413902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1305,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "926271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1305,
            "unit": "ns/op",
            "extra": "926271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "926271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "926271 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 479.8,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2506842 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 479.8,
            "unit": "ns/op",
            "extra": "2506842 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2506842 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2506842 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 859.2,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1409998 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 859.2,
            "unit": "ns/op",
            "extra": "1409998 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1409998 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1409998 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "53933f218bd4f1b501f51e3365dea01ef6dd0bf2",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-11T03:03:01Z",
          "tree_id": "3f7317f350dd976f3b5e0f391c4715096d2e19b6",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/53933f218bd4f1b501f51e3365dea01ef6dd0bf2"
        },
        "date": 1720667396914,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 426.6,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2789930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 426.6,
            "unit": "ns/op",
            "extra": "2789930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2789930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2789930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8570,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "122420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8570,
            "unit": "ns/op",
            "extra": "122420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "122420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "122420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5723,
            "unit": "ns/op\t     640 B/op\t       8 allocs/op",
            "extra": "184880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5723,
            "unit": "ns/op",
            "extra": "184880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "184880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "184880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9622,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "122908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9622,
            "unit": "ns/op",
            "extra": "122908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "122908 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "122908 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12605,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "101184 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12605,
            "unit": "ns/op",
            "extra": "101184 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 240,
            "unit": "B/op",
            "extra": "101184 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "101184 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 689.9,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1695020 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 689.9,
            "unit": "ns/op",
            "extra": "1695020 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1695020 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1695020 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5263037,
            "unit": "ns/op\t  815868 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5263037,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815868,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000772,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000772,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13719,
            "unit": "ns/op\t     355 B/op\t      10 allocs/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13719,
            "unit": "ns/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "96484 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 859.2,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1401129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 859.2,
            "unit": "ns/op",
            "extra": "1401129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1401129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1401129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13571,
            "unit": "ns/op\t     544 B/op\t      13 allocs/op",
            "extra": "85237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13571,
            "unit": "ns/op",
            "extra": "85237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "85237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "85237 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6296,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6296,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.624,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.624,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6217,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6186,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6186,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 871.7,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1410585 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 871.7,
            "unit": "ns/op",
            "extra": "1410585 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1410585 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1410585 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 861.6,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1412133 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 861.6,
            "unit": "ns/op",
            "extra": "1412133 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1412133 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1412133 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 865.6,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1409115 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 865.6,
            "unit": "ns/op",
            "extra": "1409115 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1409115 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1409115 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 985.9,
            "unit": "ns/op\t     390 B/op\t       4 allocs/op",
            "extra": "1265200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 985.9,
            "unit": "ns/op",
            "extra": "1265200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 390,
            "unit": "B/op",
            "extra": "1265200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1265200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 882.3,
            "unit": "ns/op\t     363 B/op\t       4 allocs/op",
            "extra": "1380129 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 882.3,
            "unit": "ns/op",
            "extra": "1380129 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 363,
            "unit": "B/op",
            "extra": "1380129 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1380129 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1287,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "901887 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1287,
            "unit": "ns/op",
            "extra": "901887 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "901887 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "901887 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 499.9,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2390724 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 499.9,
            "unit": "ns/op",
            "extra": "2390724 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2390724 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2390724 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 841.6,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1434216 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 841.6,
            "unit": "ns/op",
            "extra": "1434216 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1434216 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1434216 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "4a2ce95dfa1f309d0c18f121b639d6effba94bd5",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-13T03:01:44Z",
          "tree_id": "9666074b132d45643dd02c91467f80b4773faa59",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/4a2ce95dfa1f309d0c18f121b639d6effba94bd5"
        },
        "date": 1720840090313,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 428.9,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2746755 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 428.9,
            "unit": "ns/op",
            "extra": "2746755 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2746755 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2746755 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8828,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "126488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8828,
            "unit": "ns/op",
            "extra": "126488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "126488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "126488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5823,
            "unit": "ns/op\t     668 B/op\t       8 allocs/op",
            "extra": "176968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5823,
            "unit": "ns/op",
            "extra": "176968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 668,
            "unit": "B/op",
            "extra": "176968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "176968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9723,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "119776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9723,
            "unit": "ns/op",
            "extra": "119776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "119776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "119776 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13076,
            "unit": "ns/op\t     209 B/op\t       5 allocs/op",
            "extra": "90979 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13076,
            "unit": "ns/op",
            "extra": "90979 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 209,
            "unit": "B/op",
            "extra": "90979 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "90979 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 686.1,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1720548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 686.1,
            "unit": "ns/op",
            "extra": "1720548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1720548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1720548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5246546,
            "unit": "ns/op\t  815881 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5246546,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815881,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000967,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000967,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13178,
            "unit": "ns/op\t     343 B/op\t      10 allocs/op",
            "extra": "89798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13178,
            "unit": "ns/op",
            "extra": "89798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 343,
            "unit": "B/op",
            "extra": "89798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "89798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 860.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1413142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 860.6,
            "unit": "ns/op",
            "extra": "1413142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1413142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1413142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 14214,
            "unit": "ns/op\t     562 B/op\t      13 allocs/op",
            "extra": "88902 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 14214,
            "unit": "ns/op",
            "extra": "88902 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 562,
            "unit": "B/op",
            "extra": "88902 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "88902 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6197,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6197,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6244,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6244,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6215,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6215,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6201,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 866.4,
            "unit": "ns/op\t     361 B/op\t       4 allocs/op",
            "extra": "1391251 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 866.4,
            "unit": "ns/op",
            "extra": "1391251 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 361,
            "unit": "B/op",
            "extra": "1391251 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1391251 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 874.1,
            "unit": "ns/op\t     372 B/op\t       4 allocs/op",
            "extra": "1340943 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 874.1,
            "unit": "ns/op",
            "extra": "1340943 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 372,
            "unit": "B/op",
            "extra": "1340943 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1340943 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 881.7,
            "unit": "ns/op\t     365 B/op\t       4 allocs/op",
            "extra": "1373391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 881.7,
            "unit": "ns/op",
            "extra": "1373391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 365,
            "unit": "B/op",
            "extra": "1373391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1373391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 967,
            "unit": "ns/op\t     391 B/op\t       4 allocs/op",
            "extra": "1258962 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 967,
            "unit": "ns/op",
            "extra": "1258962 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 391,
            "unit": "B/op",
            "extra": "1258962 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1258962 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 911.2,
            "unit": "ns/op\t     371 B/op\t       4 allocs/op",
            "extra": "1346432 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 911.2,
            "unit": "ns/op",
            "extra": "1346432 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 371,
            "unit": "B/op",
            "extra": "1346432 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1346432 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1293,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "902444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1293,
            "unit": "ns/op",
            "extra": "902444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "902444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "902444 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 508.7,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2387016 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 508.7,
            "unit": "ns/op",
            "extra": "2387016 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2387016 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2387016 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 842.5,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1402269 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 842.5,
            "unit": "ns/op",
            "extra": "1402269 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1402269 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1402269 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "82b0004cc65cec1743e1cda8ea40f18e08dcd215",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-16T03:01:43Z",
          "tree_id": "9a849c513cafcae835a8cc3c405a3915129eb6f1",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/82b0004cc65cec1743e1cda8ea40f18e08dcd215"
        },
        "date": 1721099329045,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 427.1,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2722111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 427.1,
            "unit": "ns/op",
            "extra": "2722111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2722111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2722111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8621,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "129397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8621,
            "unit": "ns/op",
            "extra": "129397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "129397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "129397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5721,
            "unit": "ns/op\t     665 B/op\t       8 allocs/op",
            "extra": "187892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5721,
            "unit": "ns/op",
            "extra": "187892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 665,
            "unit": "B/op",
            "extra": "187892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "187892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 10174,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "118884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 10174,
            "unit": "ns/op",
            "extra": "118884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "118884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "118884 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12976,
            "unit": "ns/op\t     201 B/op\t       5 allocs/op",
            "extra": "87580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12976,
            "unit": "ns/op",
            "extra": "87580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 201,
            "unit": "B/op",
            "extra": "87580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "87580 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 729.8,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1721994 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 729.8,
            "unit": "ns/op",
            "extra": "1721994 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1721994 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1721994 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5251148,
            "unit": "ns/op\t  815864 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5251148,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815864,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0001138,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0001138,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13118,
            "unit": "ns/op\t     360 B/op\t      10 allocs/op",
            "extra": "91496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13118,
            "unit": "ns/op",
            "extra": "91496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "91496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "91496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 856.5,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1363161 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 856.5,
            "unit": "ns/op",
            "extra": "1363161 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1363161 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1363161 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13930,
            "unit": "ns/op\t     554 B/op\t      13 allocs/op",
            "extra": "83881 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13930,
            "unit": "ns/op",
            "extra": "83881 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 554,
            "unit": "B/op",
            "extra": "83881 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "83881 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6212,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6212,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6263,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6263,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6487,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6487,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6221,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6221,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 870.5,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1436368 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 870.5,
            "unit": "ns/op",
            "extra": "1436368 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1436368 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1436368 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 860.9,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1433742 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 860.9,
            "unit": "ns/op",
            "extra": "1433742 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1433742 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1433742 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 870.1,
            "unit": "ns/op\t     355 B/op\t       4 allocs/op",
            "extra": "1421286 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 870.1,
            "unit": "ns/op",
            "extra": "1421286 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "1421286 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1421286 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 972.2,
            "unit": "ns/op\t     388 B/op\t       4 allocs/op",
            "extra": "1270695 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 972.2,
            "unit": "ns/op",
            "extra": "1270695 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 388,
            "unit": "B/op",
            "extra": "1270695 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1270695 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 882.5,
            "unit": "ns/op\t     364 B/op\t       4 allocs/op",
            "extra": "1377151 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 882.5,
            "unit": "ns/op",
            "extra": "1377151 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 364,
            "unit": "B/op",
            "extra": "1377151 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1377151 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1280,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "892978 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1280,
            "unit": "ns/op",
            "extra": "892978 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "892978 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "892978 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 499.9,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2377974 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 499.9,
            "unit": "ns/op",
            "extra": "2377974 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2377974 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2377974 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 857.9,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1401373 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 857.9,
            "unit": "ns/op",
            "extra": "1401373 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1401373 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1401373 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "f69b03d12c296722c330dc928c62ad94365a3bf2",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-17T03:01:35Z",
          "tree_id": "91477ae47351a84d4a0001e0d47e038763d302d9",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/f69b03d12c296722c330dc928c62ad94365a3bf2"
        },
        "date": 1721185710264,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 428.9,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2705053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 428.9,
            "unit": "ns/op",
            "extra": "2705053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2705053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2705053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8703,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "114973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8703,
            "unit": "ns/op",
            "extra": "114973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "114973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "114973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5685,
            "unit": "ns/op\t     665 B/op\t       8 allocs/op",
            "extra": "194210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5685,
            "unit": "ns/op",
            "extra": "194210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 665,
            "unit": "B/op",
            "extra": "194210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "194210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9596,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "122043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9596,
            "unit": "ns/op",
            "extra": "122043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "122043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "122043 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 11486,
            "unit": "ns/op\t     200 B/op\t       5 allocs/op",
            "extra": "100773 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 11486,
            "unit": "ns/op",
            "extra": "100773 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 200,
            "unit": "B/op",
            "extra": "100773 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "100773 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 688.3,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1709721 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 688.3,
            "unit": "ns/op",
            "extra": "1709721 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1709721 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1709721 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5244356,
            "unit": "ns/op\t  815879 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5244356,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815879,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000811,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000811,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 11397,
            "unit": "ns/op\t     346 B/op\t      10 allocs/op",
            "extra": "97076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 11397,
            "unit": "ns/op",
            "extra": "97076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 346,
            "unit": "B/op",
            "extra": "97076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "97076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 894.1,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1366424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 894.1,
            "unit": "ns/op",
            "extra": "1366424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1366424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1366424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12850,
            "unit": "ns/op\t     551 B/op\t      13 allocs/op",
            "extra": "99669 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12850,
            "unit": "ns/op",
            "extra": "99669 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 551,
            "unit": "B/op",
            "extra": "99669 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "99669 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.62,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6225,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6225,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6192,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6192,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6193,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6193,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 869.2,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1437748 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 869.2,
            "unit": "ns/op",
            "extra": "1437748 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1437748 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1437748 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 857.1,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1444166 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 857.1,
            "unit": "ns/op",
            "extra": "1444166 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1444166 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1444166 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 863.1,
            "unit": "ns/op\t     356 B/op\t       4 allocs/op",
            "extra": "1413949 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 863.1,
            "unit": "ns/op",
            "extra": "1413949 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 356,
            "unit": "B/op",
            "extra": "1413949 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1413949 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 956.2,
            "unit": "ns/op\t     381 B/op\t       4 allocs/op",
            "extra": "1301956 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 956.2,
            "unit": "ns/op",
            "extra": "1301956 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 381,
            "unit": "B/op",
            "extra": "1301956 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1301956 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 887.8,
            "unit": "ns/op\t     362 B/op\t       4 allocs/op",
            "extra": "1385815 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 887.8,
            "unit": "ns/op",
            "extra": "1385815 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 362,
            "unit": "B/op",
            "extra": "1385815 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1385815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1298,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "935593 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1298,
            "unit": "ns/op",
            "extra": "935593 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "935593 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "935593 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 504.2,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2386179 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 504.2,
            "unit": "ns/op",
            "extra": "2386179 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2386179 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2386179 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 853.3,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1409058 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 853.3,
            "unit": "ns/op",
            "extra": "1409058 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1409058 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1409058 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "771724bfeef554f47e77eae8dbd4c0da1e36f1e5",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-20T03:01:40Z",
          "tree_id": "d13e88e69a88bc4ca554edf5cd355532200edfc4",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/771724bfeef554f47e77eae8dbd4c0da1e36f1e5"
        },
        "date": 1721444882895,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 441.8,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2711454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 441.8,
            "unit": "ns/op",
            "extra": "2711454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2711454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2711454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8585,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "117496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8585,
            "unit": "ns/op",
            "extra": "117496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "117496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "117496 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5798,
            "unit": "ns/op\t     637 B/op\t       8 allocs/op",
            "extra": "208570 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5798,
            "unit": "ns/op",
            "extra": "208570 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 637,
            "unit": "B/op",
            "extra": "208570 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "208570 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9551,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "113313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9551,
            "unit": "ns/op",
            "extra": "113313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "113313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "113313 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12706,
            "unit": "ns/op\t     209 B/op\t       5 allocs/op",
            "extra": "90916 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12706,
            "unit": "ns/op",
            "extra": "90916 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 209,
            "unit": "B/op",
            "extra": "90916 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "90916 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 703.1,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1662548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 703.1,
            "unit": "ns/op",
            "extra": "1662548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1662548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1662548 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5247183,
            "unit": "ns/op\t  815870 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5247183,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815870,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0001117,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0001117,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13588,
            "unit": "ns/op\t     354 B/op\t      10 allocs/op",
            "extra": "88150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13588,
            "unit": "ns/op",
            "extra": "88150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "88150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "88150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 873.4,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1397654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 873.4,
            "unit": "ns/op",
            "extra": "1397654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1397654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1397654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13681,
            "unit": "ns/op\t     544 B/op\t      13 allocs/op",
            "extra": "85657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13681,
            "unit": "ns/op",
            "extra": "85657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "85657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "85657 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6649,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6649,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6232,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6232,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6254,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6254,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6185,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6185,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 909,
            "unit": "ns/op\t     354 B/op\t       4 allocs/op",
            "extra": "1427354 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 909,
            "unit": "ns/op",
            "extra": "1427354 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "1427354 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1427354 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 879.3,
            "unit": "ns/op\t     364 B/op\t       4 allocs/op",
            "extra": "1378081 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 879.3,
            "unit": "ns/op",
            "extra": "1378081 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 364,
            "unit": "B/op",
            "extra": "1378081 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1378081 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 893.2,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1410708 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 893.2,
            "unit": "ns/op",
            "extra": "1410708 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1410708 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1410708 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 975.1,
            "unit": "ns/op\t     391 B/op\t       4 allocs/op",
            "extra": "1260427 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 975.1,
            "unit": "ns/op",
            "extra": "1260427 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 391,
            "unit": "B/op",
            "extra": "1260427 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1260427 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 903.3,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1398009 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 903.3,
            "unit": "ns/op",
            "extra": "1398009 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1398009 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1398009 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1321,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "884290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1321,
            "unit": "ns/op",
            "extra": "884290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "884290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "884290 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 509.3,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2338126 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 509.3,
            "unit": "ns/op",
            "extra": "2338126 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2338126 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2338126 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 869.5,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1381159 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 869.5,
            "unit": "ns/op",
            "extra": "1381159 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1381159 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1381159 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "08cc0f994287368a89cf11a33b762ec97d61657e",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-23T03:01:40Z",
          "tree_id": "e9106eb507903a0f607f511177de8f5d0dd3e96f",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/08cc0f994287368a89cf11a33b762ec97d61657e"
        },
        "date": 1721704094524,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 431.4,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2789961 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 431.4,
            "unit": "ns/op",
            "extra": "2789961 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2789961 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2789961 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8677,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "129619 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8677,
            "unit": "ns/op",
            "extra": "129619 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "129619 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "129619 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5765,
            "unit": "ns/op\t     666 B/op\t       8 allocs/op",
            "extra": "185422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5765,
            "unit": "ns/op",
            "extra": "185422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 666,
            "unit": "B/op",
            "extra": "185422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "185422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9710,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "121758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9710,
            "unit": "ns/op",
            "extra": "121758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "121758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "121758 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13143,
            "unit": "ns/op\t     210 B/op\t       5 allocs/op",
            "extra": "90640 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13143,
            "unit": "ns/op",
            "extra": "90640 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 210,
            "unit": "B/op",
            "extra": "90640 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "90640 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 704.6,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1669124 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 704.6,
            "unit": "ns/op",
            "extra": "1669124 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1669124 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1669124 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5246804,
            "unit": "ns/op\t  815874 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5246804,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815874,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0001148,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0001148,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13532,
            "unit": "ns/op\t     354 B/op\t      10 allocs/op",
            "extra": "88684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13532,
            "unit": "ns/op",
            "extra": "88684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "88684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "88684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 862.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1401141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 862.6,
            "unit": "ns/op",
            "extra": "1401141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1401141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1401141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13903,
            "unit": "ns/op\t     554 B/op\t      13 allocs/op",
            "extra": "85977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13903,
            "unit": "ns/op",
            "extra": "85977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 554,
            "unit": "B/op",
            "extra": "85977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "85977 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6195,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6195,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6331,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6331,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6193,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6193,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6196,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6196,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 863.2,
            "unit": "ns/op\t     350 B/op\t       4 allocs/op",
            "extra": "1448006 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 863.2,
            "unit": "ns/op",
            "extra": "1448006 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 350,
            "unit": "B/op",
            "extra": "1448006 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1448006 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 862.5,
            "unit": "ns/op\t     364 B/op\t       4 allocs/op",
            "extra": "1378214 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 862.5,
            "unit": "ns/op",
            "extra": "1378214 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 364,
            "unit": "B/op",
            "extra": "1378214 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1378214 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 872,
            "unit": "ns/op\t     359 B/op\t       4 allocs/op",
            "extra": "1402068 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 872,
            "unit": "ns/op",
            "extra": "1402068 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 359,
            "unit": "B/op",
            "extra": "1402068 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1402068 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 967.1,
            "unit": "ns/op\t     389 B/op\t       4 allocs/op",
            "extra": "1268253 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 967.1,
            "unit": "ns/op",
            "extra": "1268253 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 389,
            "unit": "B/op",
            "extra": "1268253 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1268253 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 882.7,
            "unit": "ns/op\t     362 B/op\t       4 allocs/op",
            "extra": "1388570 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 882.7,
            "unit": "ns/op",
            "extra": "1388570 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 362,
            "unit": "B/op",
            "extra": "1388570 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1388570 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1313,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "908938 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1313,
            "unit": "ns/op",
            "extra": "908938 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "908938 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "908938 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 518,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2379816 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 518,
            "unit": "ns/op",
            "extra": "2379816 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2379816 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2379816 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 838.5,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1431112 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 838.5,
            "unit": "ns/op",
            "extra": "1431112 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1431112 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1431112 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "8bf7a279a5b7776baa7c8bb5ca7e0bb2f22e2c32",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-07-26T03:01:55Z",
          "tree_id": "61b7fb0c357b8ab447fed7caaec7fed06eca79d4",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/8bf7a279a5b7776baa7c8bb5ca7e0bb2f22e2c32"
        },
        "date": 1721963319721,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 424.3,
            "unit": "ns/op\t     562 B/op\t       2 allocs/op",
            "extra": "2807446 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 424.3,
            "unit": "ns/op",
            "extra": "2807446 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 562,
            "unit": "B/op",
            "extra": "2807446 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2807446 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8687,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "123289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8687,
            "unit": "ns/op",
            "extra": "123289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "123289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "123289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5746,
            "unit": "ns/op\t     643 B/op\t       8 allocs/op",
            "extra": "188142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5746,
            "unit": "ns/op",
            "extra": "188142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 643,
            "unit": "B/op",
            "extra": "188142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "188142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9638,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "120909 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9638,
            "unit": "ns/op",
            "extra": "120909 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "120909 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "120909 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 11507,
            "unit": "ns/op\t     200 B/op\t       5 allocs/op",
            "extra": "100495 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 11507,
            "unit": "ns/op",
            "extra": "100495 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 200,
            "unit": "B/op",
            "extra": "100495 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "100495 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 714.1,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1682804 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 714.1,
            "unit": "ns/op",
            "extra": "1682804 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1682804 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1682804 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5237120,
            "unit": "ns/op\t  815866 B/op\t      36 allocs/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5237120,
            "unit": "ns/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815866,
            "unit": "B/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000716,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000716,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12053,
            "unit": "ns/op\t     347 B/op\t      10 allocs/op",
            "extra": "95989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12053,
            "unit": "ns/op",
            "extra": "95989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 347,
            "unit": "B/op",
            "extra": "95989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "95989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 857.7,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1374200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 857.7,
            "unit": "ns/op",
            "extra": "1374200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1374200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1374200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12555,
            "unit": "ns/op\t     559 B/op\t      13 allocs/op",
            "extra": "100290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12555,
            "unit": "ns/op",
            "extra": "100290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 559,
            "unit": "B/op",
            "extra": "100290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "100290 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6212,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6212,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6228,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6228,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6217,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6186,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6186,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 868.2,
            "unit": "ns/op\t     354 B/op\t       4 allocs/op",
            "extra": "1424282 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 868.2,
            "unit": "ns/op",
            "extra": "1424282 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "1424282 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1424282 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 853.6,
            "unit": "ns/op\t     359 B/op\t       4 allocs/op",
            "extra": "1399796 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 853.6,
            "unit": "ns/op",
            "extra": "1399796 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 359,
            "unit": "B/op",
            "extra": "1399796 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1399796 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 890.8,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1397036 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 890.8,
            "unit": "ns/op",
            "extra": "1397036 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1397036 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1397036 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 963.2,
            "unit": "ns/op\t     390 B/op\t       4 allocs/op",
            "extra": "1262578 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 963.2,
            "unit": "ns/op",
            "extra": "1262578 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 390,
            "unit": "B/op",
            "extra": "1262578 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1262578 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 882.4,
            "unit": "ns/op\t     361 B/op\t       4 allocs/op",
            "extra": "1392151 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 882.4,
            "unit": "ns/op",
            "extra": "1392151 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 361,
            "unit": "B/op",
            "extra": "1392151 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1392151 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1312,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "936979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1312,
            "unit": "ns/op",
            "extra": "936979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "936979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "936979 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 502.7,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2379736 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 502.7,
            "unit": "ns/op",
            "extra": "2379736 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2379736 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2379736 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 880.5,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1413186 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 880.5,
            "unit": "ns/op",
            "extra": "1413186 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1413186 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1413186 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "50fde94e1399dcc74cb97db32054812fcbc498bb",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-01T03:01:44Z",
          "tree_id": "6f643c17d5f78c8216b5c9884afebfd57718b3fa",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/50fde94e1399dcc74cb97db32054812fcbc498bb"
        },
        "date": 1722481702787,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 431.8,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2819708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 431.8,
            "unit": "ns/op",
            "extra": "2819708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2819708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2819708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8666,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "116151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8666,
            "unit": "ns/op",
            "extra": "116151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "116151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "116151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5758,
            "unit": "ns/op\t     629 B/op\t       8 allocs/op",
            "extra": "181314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5758,
            "unit": "ns/op",
            "extra": "181314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 629,
            "unit": "B/op",
            "extra": "181314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "181314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9823,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "119980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9823,
            "unit": "ns/op",
            "extra": "119980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "119980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "119980 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13515,
            "unit": "ns/op\t     239 B/op\t       5 allocs/op",
            "extra": "86220 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13515,
            "unit": "ns/op",
            "extra": "86220 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 239,
            "unit": "B/op",
            "extra": "86220 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "86220 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 702.1,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1702742 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 702.1,
            "unit": "ns/op",
            "extra": "1702742 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1702742 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1702742 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5267793,
            "unit": "ns/op\t  815880 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5267793,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815880,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000866,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000866,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 14048,
            "unit": "ns/op\t     351 B/op\t      10 allocs/op",
            "extra": "91826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 14048,
            "unit": "ns/op",
            "extra": "91826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 351,
            "unit": "B/op",
            "extra": "91826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "91826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 842.9,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1425696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 842.9,
            "unit": "ns/op",
            "extra": "1425696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1425696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1425696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12506,
            "unit": "ns/op\t     561 B/op\t      13 allocs/op",
            "extra": "93159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12506,
            "unit": "ns/op",
            "extra": "93159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "93159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "93159 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6217,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6313,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6313,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6207,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6211,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6211,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 872.8,
            "unit": "ns/op\t     354 B/op\t       4 allocs/op",
            "extra": "1426818 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 872.8,
            "unit": "ns/op",
            "extra": "1426818 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "1426818 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1426818 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 862.8,
            "unit": "ns/op\t     363 B/op\t       4 allocs/op",
            "extra": "1383200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 862.8,
            "unit": "ns/op",
            "extra": "1383200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 363,
            "unit": "B/op",
            "extra": "1383200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1383200 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 868.2,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1394976 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 868.2,
            "unit": "ns/op",
            "extra": "1394976 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1394976 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1394976 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 983.3,
            "unit": "ns/op\t     386 B/op\t       4 allocs/op",
            "extra": "1279687 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 983.3,
            "unit": "ns/op",
            "extra": "1279687 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 386,
            "unit": "B/op",
            "extra": "1279687 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1279687 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 883.6,
            "unit": "ns/op\t     364 B/op\t       4 allocs/op",
            "extra": "1375755 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 883.6,
            "unit": "ns/op",
            "extra": "1375755 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 364,
            "unit": "B/op",
            "extra": "1375755 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1375755 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1301,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "877567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1301,
            "unit": "ns/op",
            "extra": "877567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "877567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "877567 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 502.4,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2386273 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 502.4,
            "unit": "ns/op",
            "extra": "2386273 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2386273 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2386273 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 846.8,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1401567 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 846.8,
            "unit": "ns/op",
            "extra": "1401567 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1401567 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1401567 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "71589f93f13137db5f5f8228aa98d5f83855590b",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-05T03:01:41Z",
          "tree_id": "baa8c8876dfa2250160e15447741ea2d46280477",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/71589f93f13137db5f5f8228aa98d5f83855590b"
        },
        "date": 1722827276883,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 431.4,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2821594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 431.4,
            "unit": "ns/op",
            "extra": "2821594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2821594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2821594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8643,
            "unit": "ns/op\t     501 B/op\t      23 allocs/op",
            "extra": "127930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8643,
            "unit": "ns/op",
            "extra": "127930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 501,
            "unit": "B/op",
            "extra": "127930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "127930 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5764,
            "unit": "ns/op\t     639 B/op\t       8 allocs/op",
            "extra": "188096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5764,
            "unit": "ns/op",
            "extra": "188096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 639,
            "unit": "B/op",
            "extra": "188096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "188096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9750,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "120739 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9750,
            "unit": "ns/op",
            "extra": "120739 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "120739 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "120739 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 11796,
            "unit": "ns/op\t     200 B/op\t       5 allocs/op",
            "extra": "99452 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 11796,
            "unit": "ns/op",
            "extra": "99452 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 200,
            "unit": "B/op",
            "extra": "99452 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "99452 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 722.1,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1737290 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 722.1,
            "unit": "ns/op",
            "extra": "1737290 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1737290 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1737290 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5234416,
            "unit": "ns/op\t  815866 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5234416,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815866,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000929,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000929,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12730,
            "unit": "ns/op\t     356 B/op\t      10 allocs/op",
            "extra": "95223 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12730,
            "unit": "ns/op",
            "extra": "95223 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 356,
            "unit": "B/op",
            "extra": "95223 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "95223 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 865.3,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1397642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 865.3,
            "unit": "ns/op",
            "extra": "1397642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1397642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1397642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12979,
            "unit": "ns/op\t     552 B/op\t      13 allocs/op",
            "extra": "92272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12979,
            "unit": "ns/op",
            "extra": "92272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "92272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "92272 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6198,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6198,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6211,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6211,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6211,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6211,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6202,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6202,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 873.8,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1395260 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 873.8,
            "unit": "ns/op",
            "extra": "1395260 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1395260 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1395260 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 871.3,
            "unit": "ns/op\t     361 B/op\t       4 allocs/op",
            "extra": "1392105 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 871.3,
            "unit": "ns/op",
            "extra": "1392105 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 361,
            "unit": "B/op",
            "extra": "1392105 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1392105 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 876.1,
            "unit": "ns/op\t     358 B/op\t       4 allocs/op",
            "extra": "1406726 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 876.1,
            "unit": "ns/op",
            "extra": "1406726 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 358,
            "unit": "B/op",
            "extra": "1406726 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1406726 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 974,
            "unit": "ns/op\t     392 B/op\t       4 allocs/op",
            "extra": "1256382 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 974,
            "unit": "ns/op",
            "extra": "1256382 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 392,
            "unit": "B/op",
            "extra": "1256382 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1256382 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 895.4,
            "unit": "ns/op\t     370 B/op\t       4 allocs/op",
            "extra": "1349101 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 895.4,
            "unit": "ns/op",
            "extra": "1349101 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 370,
            "unit": "B/op",
            "extra": "1349101 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1349101 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1296,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "895429 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1296,
            "unit": "ns/op",
            "extra": "895429 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "895429 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "895429 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 506.7,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2369328 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 506.7,
            "unit": "ns/op",
            "extra": "2369328 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2369328 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2369328 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 860.4,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1403635 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 860.4,
            "unit": "ns/op",
            "extra": "1403635 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1403635 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1403635 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "3eed8b24c4245e90a0324d3a91e16064f350f8f0",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-06T03:01:47Z",
          "tree_id": "f9632a5a8c74d96819f987adfa3d6e32ec571d11",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/3eed8b24c4245e90a0324d3a91e16064f350f8f0"
        },
        "date": 1722913734820,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 428.3,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2750518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 428.3,
            "unit": "ns/op",
            "extra": "2750518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2750518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2750518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8646,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "123290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8646,
            "unit": "ns/op",
            "extra": "123290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "123290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "123290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5760,
            "unit": "ns/op\t     683 B/op\t       8 allocs/op",
            "extra": "187314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5760,
            "unit": "ns/op",
            "extra": "187314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 683,
            "unit": "B/op",
            "extra": "187314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "187314 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9765,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "120244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9765,
            "unit": "ns/op",
            "extra": "120244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "120244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "120244 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12617,
            "unit": "ns/op\t     216 B/op\t       5 allocs/op",
            "extra": "98300 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12617,
            "unit": "ns/op",
            "extra": "98300 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 216,
            "unit": "B/op",
            "extra": "98300 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "98300 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 695.2,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1636012 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 695.2,
            "unit": "ns/op",
            "extra": "1636012 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1636012 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1636012 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5252627,
            "unit": "ns/op\t  815869 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5252627,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815869,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000795,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000795,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 14041,
            "unit": "ns/op\t     346 B/op\t      10 allocs/op",
            "extra": "86727 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 14041,
            "unit": "ns/op",
            "extra": "86727 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 346,
            "unit": "B/op",
            "extra": "86727 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "86727 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 846.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1413196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 846.6,
            "unit": "ns/op",
            "extra": "1413196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1413196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1413196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 14077,
            "unit": "ns/op\t     553 B/op\t      13 allocs/op",
            "extra": "91054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 14077,
            "unit": "ns/op",
            "extra": "91054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 553,
            "unit": "B/op",
            "extra": "91054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "91054 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6217,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6283,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6283,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6219,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6219,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6203,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6203,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 861.5,
            "unit": "ns/op\t     354 B/op\t       4 allocs/op",
            "extra": "1423348 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 861.5,
            "unit": "ns/op",
            "extra": "1423348 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "1423348 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1423348 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 870.5,
            "unit": "ns/op\t     358 B/op\t       4 allocs/op",
            "extra": "1405004 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 870.5,
            "unit": "ns/op",
            "extra": "1405004 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 358,
            "unit": "B/op",
            "extra": "1405004 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1405004 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 878.5,
            "unit": "ns/op\t     359 B/op\t       4 allocs/op",
            "extra": "1398903 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 878.5,
            "unit": "ns/op",
            "extra": "1398903 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 359,
            "unit": "B/op",
            "extra": "1398903 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1398903 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 994.8,
            "unit": "ns/op\t     395 B/op\t       4 allocs/op",
            "extra": "1245206 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 994.8,
            "unit": "ns/op",
            "extra": "1245206 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 395,
            "unit": "B/op",
            "extra": "1245206 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1245206 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 892.6,
            "unit": "ns/op\t     370 B/op\t       4 allocs/op",
            "extra": "1350919 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 892.6,
            "unit": "ns/op",
            "extra": "1350919 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 370,
            "unit": "B/op",
            "extra": "1350919 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1350919 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1292,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "904081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1292,
            "unit": "ns/op",
            "extra": "904081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "904081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "904081 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 498.7,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2382817 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 498.7,
            "unit": "ns/op",
            "extra": "2382817 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2382817 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2382817 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 841.1,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1412354 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 841.1,
            "unit": "ns/op",
            "extra": "1412354 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1412354 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1412354 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "3e0dffb89839eff69a74f8f8a15e553393f64266",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-07T03:01:44Z",
          "tree_id": "0c8ee4565565a085b75e311a31dce5a4eb5e95f5",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/3e0dffb89839eff69a74f8f8a15e553393f64266"
        },
        "date": 1723000091478,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 433.2,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2794351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 433.2,
            "unit": "ns/op",
            "extra": "2794351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2794351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2794351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8700,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "123334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8700,
            "unit": "ns/op",
            "extra": "123334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "123334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "123334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5797,
            "unit": "ns/op\t     630 B/op\t       8 allocs/op",
            "extra": "179943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5797,
            "unit": "ns/op",
            "extra": "179943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 630,
            "unit": "B/op",
            "extra": "179943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "179943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9681,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "120177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9681,
            "unit": "ns/op",
            "extra": "120177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "120177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "120177 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13128,
            "unit": "ns/op\t     228 B/op\t       5 allocs/op",
            "extra": "90594 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13128,
            "unit": "ns/op",
            "extra": "90594 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 228,
            "unit": "B/op",
            "extra": "90594 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "90594 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 705.7,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1650031 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 705.7,
            "unit": "ns/op",
            "extra": "1650031 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1650031 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1650031 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5256995,
            "unit": "ns/op\t  815890 B/op\t      36 allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5256995,
            "unit": "ns/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815890,
            "unit": "B/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000955,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000955,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13577,
            "unit": "ns/op\t     365 B/op\t      10 allocs/op",
            "extra": "87159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13577,
            "unit": "ns/op",
            "extra": "87159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 365,
            "unit": "B/op",
            "extra": "87159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "87159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 868.7,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1397184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 868.7,
            "unit": "ns/op",
            "extra": "1397184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1397184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1397184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13993,
            "unit": "ns/op\t     544 B/op\t      13 allocs/op",
            "extra": "86682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13993,
            "unit": "ns/op",
            "extra": "86682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "86682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "86682 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6256,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6256,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6316,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6316,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6202,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6202,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6192,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6192,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 879.9,
            "unit": "ns/op\t     353 B/op\t       4 allocs/op",
            "extra": "1429022 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 879.9,
            "unit": "ns/op",
            "extra": "1429022 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 353,
            "unit": "B/op",
            "extra": "1429022 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1429022 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 867.6,
            "unit": "ns/op\t     365 B/op\t       4 allocs/op",
            "extra": "1373178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 867.6,
            "unit": "ns/op",
            "extra": "1373178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 365,
            "unit": "B/op",
            "extra": "1373178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1373178 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 883.9,
            "unit": "ns/op\t     356 B/op\t       4 allocs/op",
            "extra": "1417606 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 883.9,
            "unit": "ns/op",
            "extra": "1417606 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 356,
            "unit": "B/op",
            "extra": "1417606 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1417606 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 982.8,
            "unit": "ns/op\t     385 B/op\t       4 allocs/op",
            "extra": "1286202 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 982.8,
            "unit": "ns/op",
            "extra": "1286202 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 385,
            "unit": "B/op",
            "extra": "1286202 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1286202 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 890.6,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1397845 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 890.6,
            "unit": "ns/op",
            "extra": "1397845 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1397845 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1397845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1299,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "817051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1299,
            "unit": "ns/op",
            "extra": "817051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "817051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "817051 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 503.1,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2376494 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 503.1,
            "unit": "ns/op",
            "extra": "2376494 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2376494 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2376494 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 846.9,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1417179 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 846.9,
            "unit": "ns/op",
            "extra": "1417179 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1417179 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1417179 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "75270008dc784fa5540f9f126dde34f066bb134b",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-08T03:01:26Z",
          "tree_id": "e9128b4d9e70452a83f45ef03749d4fcc7cc1f3e",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/75270008dc784fa5540f9f126dde34f066bb134b"
        },
        "date": 1723086487597,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 432.3,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2764454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 432.3,
            "unit": "ns/op",
            "extra": "2764454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2764454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2764454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8986,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "119138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8986,
            "unit": "ns/op",
            "extra": "119138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "119138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "119138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5787,
            "unit": "ns/op\t     639 B/op\t       8 allocs/op",
            "extra": "179038 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5787,
            "unit": "ns/op",
            "extra": "179038 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 639,
            "unit": "B/op",
            "extra": "179038 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "179038 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9846,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "117566 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9846,
            "unit": "ns/op",
            "extra": "117566 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "117566 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "117566 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 11664,
            "unit": "ns/op\t     225 B/op\t       5 allocs/op",
            "extra": "98484 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 11664,
            "unit": "ns/op",
            "extra": "98484 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 225,
            "unit": "B/op",
            "extra": "98484 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "98484 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 722.6,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1727323 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 722.6,
            "unit": "ns/op",
            "extra": "1727323 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1727323 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1727323 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5237457,
            "unit": "ns/op\t  815876 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5237457,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815876,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000561,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000561,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12687,
            "unit": "ns/op\t     358 B/op\t      10 allocs/op",
            "extra": "93657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12687,
            "unit": "ns/op",
            "extra": "93657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 358,
            "unit": "B/op",
            "extra": "93657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "93657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 869.4,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1375136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 869.4,
            "unit": "ns/op",
            "extra": "1375136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1375136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1375136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12883,
            "unit": "ns/op\t     561 B/op\t      13 allocs/op",
            "extra": "91951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12883,
            "unit": "ns/op",
            "extra": "91951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "91951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "91951 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6234,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6234,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6208,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6209,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6209,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6235,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6235,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 872.5,
            "unit": "ns/op\t     354 B/op\t       4 allocs/op",
            "extra": "1425390 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 872.5,
            "unit": "ns/op",
            "extra": "1425390 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "1425390 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1425390 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 865.8,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1433419 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 865.8,
            "unit": "ns/op",
            "extra": "1433419 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1433419 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1433419 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 869.1,
            "unit": "ns/op\t     362 B/op\t       4 allocs/op",
            "extra": "1386391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 869.1,
            "unit": "ns/op",
            "extra": "1386391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 362,
            "unit": "B/op",
            "extra": "1386391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1386391 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 957.9,
            "unit": "ns/op\t     384 B/op\t       4 allocs/op",
            "extra": "1287555 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 957.9,
            "unit": "ns/op",
            "extra": "1287555 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "1287555 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1287555 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 871.6,
            "unit": "ns/op\t     355 B/op\t       4 allocs/op",
            "extra": "1421568 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 871.6,
            "unit": "ns/op",
            "extra": "1421568 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "1421568 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1421568 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1308,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "902054 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1308,
            "unit": "ns/op",
            "extra": "902054 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "902054 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "902054 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 505.4,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2366223 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 505.4,
            "unit": "ns/op",
            "extra": "2366223 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2366223 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2366223 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 856.7,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1405502 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 856.7,
            "unit": "ns/op",
            "extra": "1405502 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1405502 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1405502 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "7e503a70fd957903715cd4ac3600dcc0e24878ab",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-09T03:01:42Z",
          "tree_id": "56f8425df39ad01bec02139e2a33d5d6844b9d21",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/7e503a70fd957903715cd4ac3600dcc0e24878ab"
        },
        "date": 1723172899227,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 470,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2371380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 470,
            "unit": "ns/op",
            "extra": "2371380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2371380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2371380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8534,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "118368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8534,
            "unit": "ns/op",
            "extra": "118368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "118368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "118368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5777,
            "unit": "ns/op\t     621 B/op\t       8 allocs/op",
            "extra": "188733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5777,
            "unit": "ns/op",
            "extra": "188733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 621,
            "unit": "B/op",
            "extra": "188733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "188733 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9697,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "123765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9697,
            "unit": "ns/op",
            "extra": "123765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "123765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "123765 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 11790,
            "unit": "ns/op\t     232 B/op\t       5 allocs/op",
            "extra": "100058 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 11790,
            "unit": "ns/op",
            "extra": "100058 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 232,
            "unit": "B/op",
            "extra": "100058 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "100058 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 690.7,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1725392 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 690.7,
            "unit": "ns/op",
            "extra": "1725392 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1725392 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1725392 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5256940,
            "unit": "ns/op\t  815888 B/op\t      36 allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5256940,
            "unit": "ns/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815888,
            "unit": "B/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "228 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000526,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000526,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 12738,
            "unit": "ns/op\t     354 B/op\t      10 allocs/op",
            "extra": "97587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 12738,
            "unit": "ns/op",
            "extra": "97587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "97587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "97587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 878.7,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1405380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 878.7,
            "unit": "ns/op",
            "extra": "1405380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1405380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1405380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 12299,
            "unit": "ns/op\t     561 B/op\t      13 allocs/op",
            "extra": "92212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 12299,
            "unit": "ns/op",
            "extra": "92212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "92212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "92212 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6195,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6195,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6207,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6208,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.619,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.619,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 853.2,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1438000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 853.2,
            "unit": "ns/op",
            "extra": "1438000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1438000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1438000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 858.8,
            "unit": "ns/op\t     354 B/op\t       4 allocs/op",
            "extra": "1425934 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 858.8,
            "unit": "ns/op",
            "extra": "1425934 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 354,
            "unit": "B/op",
            "extra": "1425934 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1425934 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 846.4,
            "unit": "ns/op\t     357 B/op\t       4 allocs/op",
            "extra": "1408848 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 846.4,
            "unit": "ns/op",
            "extra": "1408848 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 357,
            "unit": "B/op",
            "extra": "1408848 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1408848 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 976.3,
            "unit": "ns/op\t     386 B/op\t       4 allocs/op",
            "extra": "1281454 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 976.3,
            "unit": "ns/op",
            "extra": "1281454 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 386,
            "unit": "B/op",
            "extra": "1281454 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1281454 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 867.3,
            "unit": "ns/op\t     364 B/op\t       4 allocs/op",
            "extra": "1378137 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 867.3,
            "unit": "ns/op",
            "extra": "1378137 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 364,
            "unit": "B/op",
            "extra": "1378137 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1378137 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1273,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "942471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1273,
            "unit": "ns/op",
            "extra": "942471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "942471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "942471 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 504.1,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2392804 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 504.1,
            "unit": "ns/op",
            "extra": "2392804 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2392804 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2392804 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 848.1,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1374134 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 848.1,
            "unit": "ns/op",
            "extra": "1374134 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1374134 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1374134 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukaszraczylo@users.noreply.github.com",
            "name": "lukaszraczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "distinct": true,
          "id": "016374722d2a6980144a7d8719019def9e6ca1a6",
          "message": "Update go.mod and go.sum\n\nSigned-off-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2024-08-10T03:01:52Z",
          "tree_id": "e56dfdb91dc269869d280850548e9b71538d3bbc",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/016374722d2a6980144a7d8719019def9e6ca1a6"
        },
        "date": 1723259297292,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 436.6,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2751208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 436.6,
            "unit": "ns/op",
            "extra": "2751208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2751208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2751208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8694,
            "unit": "ns/op\t     502 B/op\t      23 allocs/op",
            "extra": "123600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8694,
            "unit": "ns/op",
            "extra": "123600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 502,
            "unit": "B/op",
            "extra": "123600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "123600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5764,
            "unit": "ns/op\t     629 B/op\t       8 allocs/op",
            "extra": "181720 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5764,
            "unit": "ns/op",
            "extra": "181720 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 629,
            "unit": "B/op",
            "extra": "181720 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "181720 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9740,
            "unit": "ns/op\t    1207 B/op\t      38 allocs/op",
            "extra": "107371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9740,
            "unit": "ns/op",
            "extra": "107371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1207,
            "unit": "B/op",
            "extra": "107371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "107371 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 12181,
            "unit": "ns/op\t     210 B/op\t       5 allocs/op",
            "extra": "88969 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 12181,
            "unit": "ns/op",
            "extra": "88969 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 210,
            "unit": "B/op",
            "extra": "88969 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "88969 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 694,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1719519 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 694,
            "unit": "ns/op",
            "extra": "1719519 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1719519 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1719519 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5255251,
            "unit": "ns/op\t  815879 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5255251,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815879,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000683,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000683,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 13430,
            "unit": "ns/op\t     377 B/op\t      10 allocs/op",
            "extra": "85156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 13430,
            "unit": "ns/op",
            "extra": "85156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 377,
            "unit": "B/op",
            "extra": "85156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "85156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 882.8,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1395182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 882.8,
            "unit": "ns/op",
            "extra": "1395182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1395182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1395182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 13001,
            "unit": "ns/op\t     561 B/op\t      13 allocs/op",
            "extra": "94333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 13001,
            "unit": "ns/op",
            "extra": "94333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "94333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "94333 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.6217,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6223,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6223,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6207,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6191,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6191,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 861.2,
            "unit": "ns/op\t     352 B/op\t       4 allocs/op",
            "extra": "1434523 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 861.2,
            "unit": "ns/op",
            "extra": "1434523 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "1434523 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1434523 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 865.5,
            "unit": "ns/op\t     355 B/op\t       4 allocs/op",
            "extra": "1419102 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 865.5,
            "unit": "ns/op",
            "extra": "1419102 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 355,
            "unit": "B/op",
            "extra": "1419102 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1419102 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 877.4,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1395440 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 877.4,
            "unit": "ns/op",
            "extra": "1395440 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1395440 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1395440 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 966.4,
            "unit": "ns/op\t     395 B/op\t       4 allocs/op",
            "extra": "1246159 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 966.4,
            "unit": "ns/op",
            "extra": "1246159 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 395,
            "unit": "B/op",
            "extra": "1246159 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1246159 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 880.8,
            "unit": "ns/op\t     359 B/op\t       4 allocs/op",
            "extra": "1402198 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 880.8,
            "unit": "ns/op",
            "extra": "1402198 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 359,
            "unit": "B/op",
            "extra": "1402198 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1402198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1319,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "922888 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1319,
            "unit": "ns/op",
            "extra": "922888 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "922888 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "922888 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 506.3,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2339131 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 506.3,
            "unit": "ns/op",
            "extra": "2339131 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2339131 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2339131 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 857.9,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1396093 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 857.9,
            "unit": "ns/op",
            "extra": "1396093 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1396093 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1396093 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "committer": {
            "email": "lukasz@raczylo.com",
            "name": "Lukasz Raczylo",
            "username": "lukaszraczylo"
          },
          "distinct": true,
          "id": "dc9e0906fd8903f9768c1324609849b0bbf011cb",
          "message": "Resolve issue when proxy could panic.\n\nIssue occured when cache was disabled via environment variables but\ngraphql queries contained the cache directive.",
          "timestamp": "2024-08-19T11:27:06+01:00",
          "tree_id": "1785db262f7b1d6f82c4de0c60f235afcc47e37a",
          "url": "https://github.com/lukaszraczylo/graphql-monitoring-proxy/commit/dc9e0906fd8903f9768c1324609849b0bbf011cb"
        },
        "date": 1724063603982,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkCacheLookupInMemory",
            "value": 438,
            "unit": "ns/op\t     563 B/op\t       2 allocs/op",
            "extra": "2756913 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - ns/op",
            "value": 438,
            "unit": "ns/op",
            "extra": "2756913 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - B/op",
            "value": 563,
            "unit": "B/op",
            "extra": "2756913 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupInMemory - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2756913 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis",
            "value": 8716,
            "unit": "ns/op\t     501 B/op\t      23 allocs/op",
            "extra": "136264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - ns/op",
            "value": 8716,
            "unit": "ns/op",
            "extra": "136264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - B/op",
            "value": 501,
            "unit": "B/op",
            "extra": "136264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheLookupRedis - allocs/op",
            "value": 23,
            "unit": "allocs/op",
            "extra": "136264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory",
            "value": 5899,
            "unit": "ns/op\t     646 B/op\t       8 allocs/op",
            "extra": "182618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - ns/op",
            "value": 5899,
            "unit": "ns/op",
            "extra": "182618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - B/op",
            "value": 646,
            "unit": "B/op",
            "extra": "182618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreInMemory - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "182618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis",
            "value": 9671,
            "unit": "ns/op\t    1206 B/op\t      38 allocs/op",
            "extra": "121588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - ns/op",
            "value": 9671,
            "unit": "ns/op",
            "extra": "121588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - B/op",
            "value": 1206,
            "unit": "B/op",
            "extra": "121588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStoreRedis - allocs/op",
            "value": 38,
            "unit": "allocs/op",
            "extra": "121588 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet",
            "value": 13682,
            "unit": "ns/op\t     220 B/op\t       5 allocs/op",
            "extra": "85662 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - ns/op",
            "value": 13682,
            "unit": "ns/op",
            "extra": "85662 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - B/op",
            "value": 220,
            "unit": "B/op",
            "extra": "85662 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheSet - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "85662 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet",
            "value": 679.1,
            "unit": "ns/op\t     561 B/op\t       2 allocs/op",
            "extra": "1765014 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - ns/op",
            "value": 679.1,
            "unit": "ns/op",
            "extra": "1765014 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - B/op",
            "value": 561,
            "unit": "B/op",
            "extra": "1765014 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1765014 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire",
            "value": 5264521,
            "unit": "ns/op\t  815873 B/op\t      36 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - ns/op",
            "value": 5264521,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - B/op",
            "value": 815873,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheExpire - allocs/op",
            "value": 36,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats",
            "value": 0.0000958,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - ns/op",
            "value": 0.0000958,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMemCacheStats - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet",
            "value": 14358,
            "unit": "ns/op\t     369 B/op\t      10 allocs/op",
            "extra": "83822 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - ns/op",
            "value": 14358,
            "unit": "ns/op",
            "extra": "83822 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - B/op",
            "value": 369,
            "unit": "B/op",
            "extra": "83822 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheSet - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "83822 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet",
            "value": 848.9,
            "unit": "ns/op\t     560 B/op\t       2 allocs/op",
            "extra": "1436410 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - ns/op",
            "value": 848.9,
            "unit": "ns/op",
            "extra": "1436410 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "1436410 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheGet - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1436410 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete",
            "value": 14226,
            "unit": "ns/op\t     553 B/op\t      13 allocs/op",
            "extra": "87405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - ns/op",
            "value": 14226,
            "unit": "ns/op",
            "extra": "87405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - B/op",
            "value": 553,
            "unit": "B/op",
            "extra": "87405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheDelete - allocs/op",
            "value": 13,
            "unit": "allocs/op",
            "extra": "87405 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew",
            "value": 0.9306,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - ns/op",
            "value": 0.9306,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNew - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat",
            "value": 0.6212,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - ns/op",
            "value": 0.6212,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormat - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel",
            "value": 0.6336,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - ns/op",
            "value": 0.6336,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel",
            "value": 0.6189,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - ns/op",
            "value": 0.6189,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_NewLogger/BenchmarkNewChangeTimeFormatAndLogLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug",
            "value": 836.5,
            "unit": "ns/op\t     342 B/op\t       4 allocs/op",
            "extra": "1490794 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - ns/op",
            "value": 836.5,
            "unit": "ns/op",
            "extra": "1490794 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - B/op",
            "value": 342,
            "unit": "B/op",
            "extra": "1490794 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Debug - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1490794 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info",
            "value": 843.3,
            "unit": "ns/op\t     345 B/op\t       4 allocs/op",
            "extra": "1473252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - ns/op",
            "value": 843.3,
            "unit": "ns/op",
            "extra": "1473252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - B/op",
            "value": 345,
            "unit": "B/op",
            "extra": "1473252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Info - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1473252 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn",
            "value": 834.4,
            "unit": "ns/op\t     348 B/op\t       4 allocs/op",
            "extra": "1456455 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - ns/op",
            "value": 834.4,
            "unit": "ns/op",
            "extra": "1456455 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - B/op",
            "value": 348,
            "unit": "B/op",
            "extra": "1456455 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Warn - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1456455 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error",
            "value": 971.8,
            "unit": "ns/op\t     383 B/op\t       4 allocs/op",
            "extra": "1294498 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - ns/op",
            "value": 971.8,
            "unit": "ns/op",
            "extra": "1294498 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - B/op",
            "value": 383,
            "unit": "B/op",
            "extra": "1294498 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Error - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1294498 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal",
            "value": 853.8,
            "unit": "ns/op\t     360 B/op\t       4 allocs/op",
            "extra": "1395752 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - ns/op",
            "value": 853.8,
            "unit": "ns/op",
            "extra": "1395752 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - B/op",
            "value": 360,
            "unit": "B/op",
            "extra": "1395752 times\n4 procs"
          },
          {
            "name": "Benchmark_Log_Fatal - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1395752 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName",
            "value": 1301,
            "unit": "ns/op\t     632 B/op\t      10 allocs/op",
            "extra": "886824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - ns/op",
            "value": 1301,
            "unit": "ns/op",
            "extra": "886824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - B/op",
            "value": 632,
            "unit": "B/op",
            "extra": "886824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetMetricsName - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "886824 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels",
            "value": 502.6,
            "unit": "ns/op\t     296 B/op\t       7 allocs/op",
            "extra": "2363665 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - ns/op",
            "value": 502.6,
            "unit": "ns/op",
            "extra": "2363665 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - B/op",
            "value": 296,
            "unit": "B/op",
            "extra": "2363665 times\n4 procs"
          },
          {
            "name": "BenchmarkCompileMetricsWithLabels - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "2363665 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName",
            "value": 858.2,
            "unit": "ns/op\t     320 B/op\t       6 allocs/op",
            "extra": "1408426 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - ns/op",
            "value": 858.2,
            "unit": "ns/op",
            "extra": "1408426 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1408426 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateMetricsName - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "1408426 times\n4 procs"
          }
        ]
      }
    ]
  }
}