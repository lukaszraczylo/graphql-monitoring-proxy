window.BENCHMARK_DATA = {
  "lastUpdate": 1719976089523,
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
      }
    ]
  }
}