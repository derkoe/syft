. test_harness.sh

DOWNLOAD_POSITIVE_CASES=0

# helper for asserting download_asset positive cases
test_helper_positive_download_asset() {
  os="$1"
  arch="$2"
  format="$3"

  # for troubleshooting
  # log_set_priority 10

  name=${PROJECT_NAME}
  github_download=$(snapshot_download_url)
  version=$(snapshot_version)

  tmpdir=$(mktemp -d)

  actual_filepath=$(download_asset "${github_download}" "${tmpdir}" "${name}" "${os}" "${arch}" "${version}" "${format}" )

  assertFileExists "${actual_filepath}" "download_asset os=${os} arch=${arch} format=${format}"

  assertFilesEqual \
    "$(snapshot_dir)/${name}_${version}_${os}_${arch}.${format}" \
    "${actual_filepath}" \
    "unable to download os=${os} arch=${arch} format=${format}"

  ((DOWNLOAD_POSITIVE_CASES++))

  rm -rf -- "$tmpdir"
}


test_download_asset_exercised_all_assets() {
  expected=$(snapshot_assets_count)

  assertEquals "${expected}" "${DOWNLOAD_POSITIVE_CASES}" "did not download all possible assets (missing an os/arch/format variant?)"
}

# helper for asserting download_asset negative cases
test_helper_negative_download_asset() {
  os="$1"
  arch="$2"
  format="$3"

  # for troubleshooting
  # log_set_priority 10

  name=${PROJECT_NAME}
  github_download=$(snapshot_download_url)
  version=$(snapshot_version)

  tmpdir=$(mktemp -d)

  # since we're expecting an error, lets not exit on errors
  set +e
  actual_filepath=$(download_asset "${github_download}" "${tmpdir}" "${name}" "${os}" "${arch}" "${version}" "${format}")

  assertNotEquals "0" "$?" "download did not indicate a failure"
  set -e

  assertEquals ""  "${actual_filepath}" "unable to download os=${os} arch=${arch} format=${format}"

  rm -rf -- "$tmpdir"
}


worker_pid=$(setup_snapshot_server)
trap 'teardown_snapshot_server ${worker_pid}' EXIT

# exercise all possible assets
run_test_case test_helper_positive_download_asset "linux" "amd64" "tar.gz"
run_test_case test_helper_positive_download_asset "linux" "amd64" "rpm"
run_test_case test_helper_positive_download_asset "linux" "amd64" "deb"
run_test_case test_helper_positive_download_asset "linux" "arm64" "tar.gz"
run_test_case test_helper_positive_download_asset "linux" "arm64" "rpm"
run_test_case test_helper_positive_download_asset "linux" "arm64" "deb"
run_test_case test_helper_positive_download_asset "darwin" "amd64" "tar.gz"
run_test_case test_helper_positive_download_asset "darwin" "amd64" "zip"
run_test_case test_helper_positive_download_asset "darwin" "arm64" "tar.gz"
run_test_case test_helper_positive_download_asset "darwin" "arm64" "zip"
run_test_case test_helper_positive_download_asset "windows" "amd64" "zip"
# note: the mac signing process produces a dmg which is not part of the snapshot process (thus is not exercised here)

# let's make certain we covered all assets that were expected
run_test_case test_download_asset_exercised_all_assets

# make certain we handle missing assets alright
run_test_case test_helper_negative_download_asset "bogus" "amd64" "zip"

trap - EXIT
teardown_snapshot_server "${worker_pid}"
