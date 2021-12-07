/*
Copyright © 2021 Adedayo Adetoye (aka Dayo)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"context"
	"fmt"
	"log"

	model "github.com/adedayo/git-service-driver/pkg"
	"github.com/adedayo/git-service-driver/pkg/gitlab"
)

func main() {

	paged := gitlab.GitLabPagedSearch{
		ServiceID:  "e6584c4f-9aa4-402f-8650-52cbf496c9f4",
		PageSize:   1000,
		First:      7,
		NextCursor: "",
	}
	config := model.MakeConfigManager().GetConfig()
	v, err := config.FindService(paged.ServiceID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	projs, loc, err := gitlab.GetRepositories(context.Background(), v, &paged)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Loc: %v\nProjs: %v\n", loc, projs)
	// cmd.Execute()
}
