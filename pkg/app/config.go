//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package app

import "github.com/algotiqa/core"

//=============================================================================

type Config struct {
	core.Application
	Routes map[string]struct {
		Prefix string
		Url    string
	}
}

//=============================================================================
