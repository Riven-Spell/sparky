`sparkrun cluster status --cluster <cluster_name> --json`

Note that different (less) information is shown if a group lives on another host... so, instead, agent health checks should exclusively ask about localhost:
`sparkrun cluster status --cluster <cluster_name> --json -H localhost`

The first host in a list will always have the HTTP(S) server, so the URI can be constructed as `http://hosts[0]:port/`.

```json
{
    "groups": {
        "sparkrun_3ba7d381d8c8": {
            "meta": {
                "cluster_id": "sparkrun_3ba7d381d8c8",
                "data_parallel": 1,
                "effective_container_image": "scitrera/dgx-spark-vllm-jasl-ds4:20260609",
                "hosts": [
                    "192.168.100.1",
                    "192.168.100.2"
                ],
                "model": "deepseek-ai/DeepSeek-V4-Flash",
                "pipeline_parallel": 1,
                "port": 8000,
                "recipe": "dsv4-flash",
                "recipe_state": {
                    "_applied_overrides": {},
                    "_qualified_name_override": null,
                    "_raw": {
                        "command": "vllm serve {model} \\\n  --served-model-name 'deepseek-ai/DeepSeek-V4-Flash' \\\n  --host {host} \\\n  --port {port} \\\n  --max-model-len {max_model_len} \\\n  --max-num-seqs {max_num_seqs} \\\n  --max-num-batched-tokens {max_num_batched_tokens} \\\n  --trust-remote-code \\\n  --gpu-memory-utilization {gpu_memory_utilization} \\\n  --enable-auto-tool-choice \\\n  --tool-call-parser {tool_call_parser} \\\n  --tokenizer-mode {tokenizer_mode} \\\n  --reasoning-parser {reasoning_parser} \\\n  --kv-cache-dtype {kv_cache_dtype} \\\n  --load-format {load_format} \\\n  --speculative-config '{speculative_config}' \\\n  --compilation-config '{compilation_config}' \\\n  --enable-prefix-caching \\\n  --max-cudagraph-capture-size {max_cudagraph_capture_size} \\\n  -tp {tensor_parallel} \\\n  -pp {pipeline_parallel} \\\n  --attention-config '{attn_config}' \\\n  --enable-flashinfer-autotune\n",
                        "container": "scitrera/dgx-spark-vllm-jasl-ds4:20260609",
                        "defaults": {
                            "attn_config": "{\"backend\": \"FLASHINFER\"}",
                            "compilation_config": "{\"cudagraph_mode\":\"FULL_AND_PIECEWISE\",\"custom_ops\":[\"all\"]}",
                            "data_parallel": 1,
                            "gpu_memory_utilization": 0.85,
                            "host": "0.0.0.0",
                            "kv_cache_dtype": "fp8",
                            "load_format": "instanttensor",
                            "max_cudagraph_capture_size": 8,
                            "max_model_len": 400000,
                            "max_num_batched_tokens": 16384,
                            "max_num_seqs": 10,
                            "pipeline_parallel": 1,
                            "port": 8000,
                            "reasoning_parser": "deepseek_v4",
                            "speculative_config": "{\"method\": \"mtp\", \"num_speculative_tokens\": 2}",
                            "tensor_parallel": 2,
                            "tokenizer_mode": "deepseek_v4",
                            "tool_call_parser": "deepseek_v4"
                        },
                        "metadata": {
                            "description": "DeepSeek-V4-Flash"
                        },
                        "min_nodes": 2,
                        "model": "deepseek-ai/DeepSeek-V4-Flash",
                        "recipe_version": "2",
                        "runtime": "vllm"
                    },
                    "_serialization_version": 1,
                    "builder": "",
                    "builder_config": {},
                    "command": "vllm serve {model} \\\n  --served-model-name 'deepseek-ai/DeepSeek-V4-Flash' \\\n  --host {host} \\\n  --port {port} \\\n  --max-model-len {max_model_len} \\\n  --max-num-seqs {max_num_seqs} \\\n  --max-num-batched-tokens {max_num_batched_tokens} \\\n  --trust-remote-code \\\n  --gpu-memory-utilization {gpu_memory_utilization} \\\n  --enable-auto-tool-choice \\\n  --tool-call-parser {tool_call_parser} \\\n  --tokenizer-mode {tokenizer_mode} \\\n  --reasoning-parser {reasoning_parser} \\\n  --kv-cache-dtype {kv_cache_dtype} \\\n  --load-format {load_format} \\\n  --speculative-config '{speculative_config}' \\\n  --compilation-config '{compilation_config}' \\\n  --enable-prefix-caching \\\n  --max-cudagraph-capture-size {max_cudagraph_capture_size} \\\n  -tp {tensor_parallel} \\\n  -pp {pipeline_parallel} \\\n  --attention-config '{attn_config}' \\\n  --enable-flashinfer-autotune\n",
                    "container": "scitrera/dgx-spark-vllm-jasl-ds4:20260609",
                    "defaults": {
                        "attn_config": "{\"backend\": \"FLASHINFER\"}",
                        "compilation_config": "{\"cudagraph_mode\":\"FULL_AND_PIECEWISE\",\"custom_ops\":[\"all\"]}",
                        "data_parallel": 1,
                        "gpu_memory_utilization": 0.85,
                        "host": "0.0.0.0",
                        "kv_cache_dtype": "fp8",
                        "load_format": "instanttensor",
                        "max_cudagraph_capture_size": 8,
                        "max_model_len": 400000,
                        "max_num_batched_tokens": 16384,
                        "max_num_seqs": 10,
                        "pipeline_parallel": 1,
                        "port": 8000,
                        "reasoning_parser": "deepseek_v4",
                        "speculative_config": "{\"method\": \"mtp\", \"num_speculative_tokens\": 2}",
                        "tensor_parallel": 2,
                        "tokenizer_mode": "deepseek_v4",
                        "tool_call_parser": "deepseek_v4"
                    },
                    "description": "DeepSeek-V4-Flash",
                    "distribution_config": {
                        "containers": {
                            "enabled": true,
                            "entries": [
                                {
                                    "name": "scitrera/dgx-spark-vllm-jasl-ds4:20260609",
                                    "target": []
                                }
                            ]
                        },
                        "externally_provided": false,
                        "models": {
                            "enabled": true,
                            "entries": [
                                {
                                    "name": "deepseek-ai/DeepSeek-V4-Flash",
                                    "revision": null,
                                    "target": []
                                }
                            ]
                        }
                    },
                    "env": {},
                    "executor_config": {},
                    "maintainer": "",
                    "max_nodes": null,
                    "metadata": {
                        "description": "DeepSeek-V4-Flash",
                        "head_dim": 512,
                        "kv_dtype": "fp8",
                        "model_dtype": "fp8",
                        "num_kv_heads": 1,
                        "num_layers": 43,
                        "quant_bits": 8,
                        "quantization": "fp8"
                    },
                    "min_nodes": 2,
                    "mode": "cluster",
                    "model": "deepseek-ai/DeepSeek-V4-Flash",
                    "model_revision": null,
                    "mods": [],
                    "name": "dsv4-flash",
                    "post_commands": [],
                    "post_exec": [],
                    "pre_exec": [],
                    "recipe_version": "2",
                    "runtime": "vllm-distributed",
                    "runtime_config": {},
                    "runtime_version": "",
                    "source_path": "dsv4-flash.yaml",
                    "source_registry": null,
                    "source_registry_url": null,
                    "stop_after_post": false
                },
                "runtime": "vllm-distributed",
                "runtime_info": {
                    "cuda": "13.2",
                    "nccl": "(2, 30, 7)",
                    "python": "3.12.3",
                    "torch": "2.12.0",
                    "vllm": "0.1.dev17473+g21a54eacf.d20260609"
                },
                "tensor_parallel": 2
            },
            "containers": [
                {
                    "host": "localhost",
                    "role": "node_0",
                    "status": "Up 9 hours",
                    "image": "scitrera/dgx-spark-vllm-jasl-ds4:20260609"
                }
            ],
            "label": "dsv4-flash  (tp=2, pp=1)  [3ba7d381d8c8]"
        }
    },
    "solo_entries": [
        {
            "cluster_id": "sparkrun_122df0c159d9",
            "meta": {
                "cluster_id": "sparkrun_122df0c159d9",
                "effective_container_image": "vllm-node-tf5",
                "hosts": [
                    "192.168.100.1"
                ],
                "model": "protoLabsAI/Ornith-1.0-35B-FP8",
                "overrides": {
                    "tensor_parallel": 1
                },
                "port": 8000,
                "recipe": "ornith-crash",
                "recipe_state": {
                    "_applied_overrides": {
                        "tensor_parallel": 1
                    },
                    "_qualified_name_override": null,
                    "_raw": {
                        "cluster_only": false,
                        "command": "vllm serve protoLabsAI/Ornith-1.0-35B-FP8 \\\n  --invalid-option \\\n  --served-model-name Ornith-1.0-35B \\\n  --tensor-parallel-size {tensor_parallel} \\\n  --host {host} --port {port} \\\n  --max-model-len {max_model_len} \\\n  --gpu-memory-utilization {gpu_memory_utilization} \\\n  --max-num-seqs {max_num_seqs} \\\n  --max-num-batched-tokens {max_num_batched_tokens} \\\n  --enable-prefix-caching \\\n  --enable-auto-tool-choice --tool-call-parser qwen3_xml \\\n  --reasoning-parser qwen3 \\\n  --trust-remote-code \\\n  --load-format {load_format} \\\n  --attention-backend {attn_backend} \\\n  --enable-chunked-prefill \\\n  --enable-flashinfer-autotune \\\n  --kv-cache-dtype {kv_cache_dtype} \\\n  --speculative-config '{\"method\":\"ngram\", \"num_speculative_tokens\": 12, \"prompt_lookup_max\": 8}'\n",
                        "container": "vllm-node-tf5",
                        "defaults": {
                            "attn_backend": "flashinfer",
                            "gpu_memory_utilization": 0.85,
                            "host": "0.0.0.0",
                            "kv_cache_dtype": "fp8",
                            "load_format": "instanttensor",
                            "max_model_len": 262144,
                            "max_num_batched_tokens": 32768,
                            "max_num_seqs": 16,
                            "port": 8000,
                            "tensor_parallel": 1
                        },
                        "description": "Ornith-1.0-35B",
                        "model": "protoLabsAI/Ornith-1.0-35B-FP8",
                        "name": "Ornith-1.0-35B",
                        "recipe_version": "1",
                        "solo_only": false
                    },
                    "_serialization_version": 1,
                    "builder": "eugr",
                    "builder_config": {},
                    "command": "vllm serve protoLabsAI/Ornith-1.0-35B-FP8 \\\n  --invalid-option \\\n  --served-model-name Ornith-1.0-35B \\\n  --tensor-parallel-size {tensor_parallel} \\\n  --host {host} --port {port} \\\n  --max-model-len {max_model_len} \\\n  --gpu-memory-utilization {gpu_memory_utilization} \\\n  --max-num-seqs {max_num_seqs} \\\n  --max-num-batched-tokens {max_num_batched_tokens} \\\n  --enable-prefix-caching \\\n  --enable-auto-tool-choice --tool-call-parser qwen3_xml \\\n  --reasoning-parser qwen3 \\\n  --trust-remote-code \\\n  --load-format {load_format} \\\n  --attention-backend {attn_backend} \\\n  --enable-chunked-prefill \\\n  --enable-flashinfer-autotune \\\n  --kv-cache-dtype {kv_cache_dtype} \\\n  --speculative-config '{\"method\":\"ngram\", \"num_speculative_tokens\": 12, \"prompt_lookup_max\": 8}'\n",
                    "container": "vllm-node-tf5",
                    "defaults": {
                        "attn_backend": "flashinfer",
                        "gpu_memory_utilization": 0.85,
                        "host": "0.0.0.0",
                        "kv_cache_dtype": "fp8",
                        "load_format": "instanttensor",
                        "max_model_len": 262144,
                        "max_num_batched_tokens": 32768,
                        "max_num_seqs": 16,
                        "port": 8000,
                        "tensor_parallel": 1
                    },
                    "description": "Ornith-1.0-35B",
                    "distribution_config": {
                        "containers": {
                            "enabled": true,
                            "entries": [
                                {
                                    "name": "vllm-node-tf5",
                                    "target": []
                                }
                            ]
                        },
                        "externally_provided": false,
                        "models": {
                            "enabled": true,
                            "entries": [
                                {
                                    "name": "protoLabsAI/Ornith-1.0-35B-FP8",
                                    "revision": null,
                                    "target": []
                                }
                            ]
                        }
                    },
                    "env": {},
                    "executor_config": {},
                    "maintainer": "",
                    "max_nodes": null,
                    "metadata": {
                        "head_dim": 256,
                        "kv_dtype": "fp8",
                        "model_dtype": "fp8",
                        "num_kv_heads": 2,
                        "num_layers": 40,
                        "quant_bits": 8,
                        "quantization": "fp8"
                    },
                    "min_nodes": 1,
                    "mode": "auto",
                    "model": "protoLabsAI/Ornith-1.0-35B-FP8",
                    "model_revision": null,
                    "mods": [],
                    "name": "ornith-crash",
                    "post_commands": [],
                    "post_exec": [],
                    "pre_exec": [],
                    "recipe_version": "1",
                    "runtime": "vllm-distributed",
                    "runtime_config": {},
                    "runtime_version": "",
                    "source_path": "ornith-crash.yaml",
                    "source_registry": null,
                    "source_registry_url": null,
                    "stop_after_post": false
                },
                "runtime": "vllm-distributed",
                "runtime_info": {
                    "build_base_image": "${CUDA_IMAGE}",
                    "build_build_args_build_jobs": "16",
                    "build_build_args_exp_mxfp4": "False",
                    "build_build_args_transformers_5": "False",
                    "build_build_args_vllm_prs": "",
                    "build_build_args_vllm_ref": "main",
                    "build_build_date": "2026-06-25 11:42:21+00:00",
                    "build_build_script_commit": "1d45b34254bc5ce4bdcdd72cf22069596c8772d2",
                    "build_flashinfer_commit": "b3baedbb",
                    "build_gpu_arch": "12.1a",
                    "build_vllm_commit": "901a3b091",
                    "build_vllm_version": "0.23.1rc1.dev309+g901a3b091.d20260623",
                    "container_maintainer": "NVIDIA CORPORATION <cudatools@nvidia.com>",
                    "container_org_opencontainers_image_ref_name": "ubuntu",
                    "container_org_opencontainers_image_version": "24.04",
                    "container_sparky_launchtime": "92570283457",
                    "cuda": "13.0",
                    "nccl": "(2, 28, 9)",
                    "python": "3.12.3",
                    "torch": "2.11.0+cu130",
                    "vllm": "0.23.1rc1.dev309+g901a3b091.d20260623"
                },
                "tensor_parallel": 1
            },
            "host": "localhost",
            "status": "Up 28 seconds",
            "image": "vllm-node-tf5",
            "label": "ornith-crash  (tp=1)  [122df0c159d9]"
        }
    ],
    "idle_hosts": [],
    "pending_ops": [],
    "errors": {},
    "total_containers": 2,
    "host_count": 1
}
```