/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"github.com/emicklei/go-restful"
)

/*
 * get the bk_biz_id from restful request, and then use the bk_biz_id
 * to get the host id list in the biz_id, and get file content of
 * /root/yovole from datacollection of every host
 */
func (s *Service) GetFileContent(req *restful.Request, resp *restful.Response) {
	// pheader := req.Request.Header
	// defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	bizID := req.PathParameter("bk_biz_id")
	blog.Debug("Yovole, bizID: %s", bizID)

	// get host id list from bk_biz_id
	// this is test data for web_server api call
	// hostList := [...]int{903}

	// get YovoleTest from datacollection
	// this is test data for web_server api call
	yovoleTest := meta.YovoleTest{BizID: 14, InnerIP: "10.3.12.210", HostID: 903, FileContent: "xyz"}

	resp.WriteEntity(meta.NewSuccessResp(yovoleTest))
}
