# SPDX-FileCopyrightText: 2025 The PolyClient Authors
#
# SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

FROM scratch
COPY polyclient /usr/bin/polyclient
ENTRYPOINT ["/usr/bin/polyclient"]