#!/bin/bash

HOST="http://localhost:8080"
INFO_ENDPOINT="/api/v1/info"
AUTH_ENDPOINT="/api/v1/auth"
SEND_COIN_ENDPOINT="/api/v1/sendCoin/"
CREDIT_ENDPOINT="/api/v1/credit"
BUY_ENDPOINT="/api/v1/buy/t-shirt"

EMPLOYEE_AUTH='{"username":"employee","password":"password"}'
DIRECTOR_AUTH='{"username":"director","password":"password"}'
SEND_COIN_DATA='{"amount":22,"toUser":"director"}'
CREDIT_DATA='{"username":"employee","amount":100000}'

TEST_DURATION="60s"
REQUESTS_PER_SECOND=1000
CONCURRENT_USERS=100
WRITE_TO_FILE=false

RESULTS_DIR="docs/load_test_results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

if [ "$WRITE_TO_FILE" = true ]; then
    mkdir -p "${RESULTS_DIR}/${TIMESTAMP}"
    echo "Results will be written to: ${RESULTS_DIR}/${TIMESTAMP}"
fi

run_oha_test() {
    local test_name=$1
    shift

    if [ "$WRITE_TO_FILE" = true ]; then
        echo "Running ${test_name} test (results will be saved to file)..."
        "$@" > "${RESULTS_DIR}/${TIMESTAMP}/${test_name}.txt" 2>&1
    else
        echo "Running ${test_name} test..."
        "$@"
    fi
}

get_token() {
    local auth_data=$1
    response=$(curl -s -L POST \
        -H "Content-Type: application/json" \
        -H "Accept: application/json" \
        -d "${auth_data}" \
        "${HOST}${AUTH_ENDPOINT}")

    if [ -z "${response}" ]; then
        echo "Error: Empty response received" >&2
        return 1
    fi

    token=$(echo "${response}" | jq -r '.token')

    if [ "$token" = "null" ] || [ -z "$token" ]; then
        echo "Error: Failed to extract token from response" >&2
        echo "Response was: ${response}" >&2
        return 1
    fi

    echo "${token}"
}

echo "Retrieving employee token..."
EMPLOYEE_TOKEN=$(get_token "${EMPLOYEE_AUTH}")

if [ -z "$EMPLOYEE_TOKEN" ]; then
    echo "Failed to retrieve employee token. Exiting."
    exit 1
fi

echo "Employee token retrieved successfully."

echo "Retrieving director token..."
DIRECTOR_TOKEN=$(get_token "${DIRECTOR_AUTH}")

if [ -z "$DIRECTOR_TOKEN" ]; then
    echo "Failed to retrieve director token. Exiting."
    exit 1
fi

echo "Director token retrieved successfully."

echo -e "\nPerforming credit operation before test..."
credit_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${DIRECTOR_TOKEN}" \
    -d "${CREDIT_DATA}" \
    "${HOST}${CREDIT_ENDPOINT}")

echo "Credit operation response: ${credit_response}"


# Auth endpoint test
run_oha_test "auth" oha -z "${TEST_DURATION}" \
    -c "${CONCURRENT_USERS}" \
    -q "${REQUESTS_PER_SECOND}" \
    --latency-correction \
    --disable-keepalive \
    -m POST \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -d "${EMPLOYEE_AUTH}" \
    "${HOST}${AUTH_ENDPOINT}"


# SendCoin endpoint test
run_oha_test "sendcoin" oha -z "${TEST_DURATION}" \
    -c "${CONCURRENT_USERS}" \
    -q "${REQUESTS_PER_SECOND}" \
    --latency-correction \
    --disable-keepalive \
    -m POST \
    -H "Accept: application/json" \
    -H "Content-Type: application/json" \
    -H "Authorization:Bearer ${EMPLOYEE_TOKEN}" \
    -d "${SEND_COIN_DATA}" \
    "${HOST}${SEND_COIN_ENDPOINT}"

# Buy endpoint test
run_oha_test "buy" oha -z "${TEST_DURATION}" \
    -c "${CONCURRENT_USERS}" \
    -q "${REQUESTS_PER_SECOND}" \
    --latency-correction \
    --disable-keepalive \
    -m POST \
    -H "Accept: application/json" \
    -H "Authorization:Bearer ${EMPLOYEE_TOKEN}" \
    "${HOST}${BUY_ENDPOINT}"

# Info endpoint test
run_oha_test "info" oha -z "${TEST_DURATION}" \
    -c "${CONCURRENT_USERS}" \
    -q "${REQUESTS_PER_SECOND}" \
    --latency-correction \
    --disable-keepalive \
    -H "Accept: application/json" \
    -H "Authorization:Bearer ${EMPLOYEE_TOKEN}" \
    "${HOST}${INFO_ENDPOINT}"

if [ "$WRITE_TO_FILE" = true ]; then
    echo -e "\nTest results have been written to: ${RESULTS_DIR}/${TIMESTAMP}"
    echo "Results files:"
    echo "- info.txt"
    echo "- auth.txt"
    echo "- sendcoin.txt"
    echo "- buy.txt"
fi
