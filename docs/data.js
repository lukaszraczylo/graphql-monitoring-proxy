window.BENCHMARK_DATA = {
  "lastUpdate": 1719594657678,
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
      }
    ]
  }
}