#!/usr/bin/env bash
set -euo pipefail

# Runs the NIST STS reference suite and compares its p-values with the Pure Go implementation.
# Defaults mimic the upstream workflow: data.pi dataset, 1,000,000 bits, all tests with default params.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BASE_DIR="${ROOT}/build/sts-2_1_2"
DATASET="data.pi"
BITCOUNT=1000000
ENCODING="auto" # auto-detect (ascii 0/1 text vs binary); override with --encoding

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dataset) DATASET="$2"; shift 2 ;;
    --bits) BITCOUNT="$2"; shift 2 ;;
    --sts-dir) BASE_DIR="$2"; shift 2 ;;
    --encoding) ENCODING="$2"; shift 2 ;;
    *) echo "Unknown arg: $1"; exit 1 ;;
  esac
done

mkdir -p "${BASE_DIR}"

# Helper to find existing assets
find_assess() { find "${BASE_DIR}" -maxdepth 3 -type f -name assess -perm -u+x | head -n1; }
find_makefile() { find "${BASE_DIR}" -maxdepth 3 -type f -name makefile | head -n1; }

ASSESS_BIN="$(find_assess)"
if [[ -z "${ASSESS_BIN}" ]]; then
  # Only download into an empty directory to avoid nested mess
  if [[ -n "$(ls -A "${BASE_DIR}")" ]]; then
    echo "Directory ${BASE_DIR} is not empty and no assess binary found. Remove it or pass --sts-dir to a clean path."
    exit 1
  fi
  echo "Downloading NIST STS..."
  curl -L -o "${BASE_DIR}/sts.zip" https://csrc.nist.rip/CSRC/media/Projects/Random-Bit-Generation/documents/sts-2_1_2.zip
  unzip -q "${BASE_DIR}/sts.zip" -d "${BASE_DIR}"
  cp -r "${BASE_DIR}/sts-2.1.2/sts-2.1.2/." ${BASE_DIR}
  rm -rf "${BASE_DIR}/sts-2.1.2"

  MAKEFILE_PATH="$(find_makefile)"
  if [[ -z "${MAKEFILE_PATH}" ]]; then
    echo "makefile not found after unzip in ${BASE_DIR}"
    exit 1
  fi
    BUILD_DIR="$(dirname "${MAKEFILE_PATH}")"
    echo "Building assess..."
    (cd "${BUILD_DIR}" && make -f makefile)
    ASSESS_BIN="${BUILD_DIR}/assess"
  else
    BUILD_DIR="$(dirname "${ASSESS_BIN}")"
    echo "Reusing existing NIST STS at ${BUILD_DIR}"
fi

DATA_PATH="${BUILD_DIR}/data/${DATASET}"
if [[ ! -f "${DATA_PATH}" ]]; then
  echo "Dataset ${DATASET} not found in ${BUILD_DIR}/data"
  exit 1
fi

if [[ "${ENCODING}" == "auto" ]]; then
  # Peek at a small sample; if anything besides 0/1/whitespace appears, treat as binary
  if head -c 4096 "${DATA_PATH}" | LC_ALL=C grep -q '[^01[:space:]]'; then
    ENCODING="binary"
  else
    ENCODING="ascii"
  fi
  echo "Detected ${ENCODING} encoding for ${DATASET}"
fi

echo "Running NIST reference suite on ${DATASET} (${BITCOUNT} bits)..."
pushd "${BUILD_DIR}" >/dev/null
# Clean prior experiment outputs to avoid mixing results.
rm -rf experiments/AlgorithmTesting
# Ensure experiments directory structure exists (AlgorithmTesting + subfolders)
if [[ -x experiments/create-dir-script ]]; then
  sh experiments/create-dir-script 2>/dev/null || true
fi
mkdir -p experiments/AlgorithmTesting
TEST_DIRS=(Frequency BlockFrequency CumulativeSums Runs LongestRun Rank FFT NonOverlappingTemplate OverlappingTemplate Universal ApproximateEntropy RandomExcursions RandomExcursionsVariant Serial LinearComplexity)
for d in "${TEST_DIRS[@]}"; do
  mkdir -p "experiments/AlgorithmTesting/${d}"
done

# Input order: generator option, file path, apply all tests, parameter adjustments (0), bitstreams count, input mode (0=ASCII)
# Assess input mode: 0=ASCII, 1=binary
INPUT_MODE=0
if [[ "${ENCODING}" == "binary" ]]; then
  INPUT_MODE=1
fi

printf "0\ndata/${DATASET}\n1\n0\n1\n${INPUT_MODE}\n" | ./assess "${BITCOUNT}" >/tmp/nist_assess.log || true
echo "Finished NIST reference suite on ${DATASET} (${BITCOUNT} bits)..."
popd >/dev/null

echo "Comparing against Pure Go implementation..."
go run ./tools/validate_nist_go_vs_c.go \
  --dataset "${DATA_PATH}" \
  --bits "${BITCOUNT}" \
  --encoding "${ENCODING}" \
  --results "${BUILD_DIR}/experiments/AlgorithmTesting"

echo "Reference run log: /tmp/nist_assess.log"
