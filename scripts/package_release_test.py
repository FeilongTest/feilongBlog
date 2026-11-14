import hashlib
import importlib.util
from pathlib import Path
import tarfile
import tempfile
import unittest

spec = importlib.util.spec_from_file_location("package_release", Path(__file__).with_name("package-release.py"))
packager = importlib.util.module_from_spec(spec)
spec.loader.exec_module(packager)


class ReleaseArchiveTest(unittest.TestCase):
    def test_mac_sidecars_are_excluded_without_dropping_real_assets(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            release = root / "release"
            assets = release / "web/vendor/vditor/dist"
            assets.mkdir(parents=True)
            (assets / "index.js").write_text("real JavaScript")
            (assets / "._index.js").write_bytes(b"AppleDouble metadata")
            (release / ".DS_Store").write_bytes(b"Finder metadata")
            (release / "._folder").mkdir()
            (release / "._folder/payload").write_bytes(b"metadata")
            output = root / "release.tar.gz"
            packager.package_release(release, output)
            with tarfile.open(output) as archive:
                names = archive.getnames()
                self.assertIn("release/web/vendor/vditor/dist/index.js", names)
                self.assertFalse(any(part.startswith("._") or part == ".DS_Store" for name in names for part in Path(name).parts))
                self.assertEqual(archive.extractfile("release/web/vendor/vditor/dist/index.js").read(), b"real JavaScript")
            self.assertEqual(output.with_name(output.name + ".sha256").read_text().split()[0], hashlib.sha256(output.read_bytes()).hexdigest())


if __name__ == "__main__":
    unittest.main()
