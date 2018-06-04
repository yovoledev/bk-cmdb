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
 
package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/migrate_service/data"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/storage"
	"sort"
	"strconv"
	"strings"
)

// Upgrade upgrade
func Upgrade(instData storage.DI) error {
	err := upgradeAppfield(instData)
	if err != nil {
		return err
	}
	err = upgradeGlobalization(instData)
	if err != nil {
		return err
	}
	return nil
}

func upgradeGlobalization(db storage.DI) error {
	presetRows := data.AppRow()
	presetRows = append(presetRows, data.HostRow()...)
	presetRows = append(presetRows, data.ModuleRow()...)
	presetRows = append(presetRows, data.PlatRow()...)
	presetRows = append(presetRows, data.ProcRow()...)
	presetRows = append(presetRows, data.SetRow()...)
	presetRowsMap := map[string]*metadata.ObjectAttDes{}
	for _, row := range presetRows {
		presetRowsMap[row.ObjectID+"::"+row.PropertyID] = row
	}

	attdes := []metadata.ObjectAttDes{}
	err := db.GetMutilByCondition(common.BKTableNameObjAttDes, nil, map[string]interface{}{common.BKPropertyTypeField: common.FieldTypeEnum}, &attdes, "", 0, 0)
	if err != nil {
		blog.Errorf("upgradeGlobalization get attdes error: %v", err)
		return err
	}
	blog.Infof("attdes:= %#v", attdes)
	for _, curRow := range attdes {
		curOptions := validator.ParseEnumOption(curRow.Option)
		blog.Infof("curOptions:= %#v", curOptions)
		if len(curOptions) <= 0 {
			continue
		}

		expectOptions := []validator.EnumVal{}
		expectRow := presetRowsMap[curRow.ObjectID+"::"+curRow.PropertyID]
		if expectRow != nil {
			var ok bool
			expectOptions, ok = expectRow.Option.([]validator.EnumVal)
			if !ok {
				expectOptions = []validator.EnumVal{}
			}
		}
		blog.Infof("expectOptions:= %#v", curOptions)

		sort.SliceStable(curOptions, func(i, j int) bool {
			return strings.Compare(curOptions[i].Name, curOptions[j].Name) < 0
		})
		// get max id
		newID := len(expectOptions)
		for _, option := range curOptions {
			id, _ := strconv.Atoi(option.ID) // ignore err cause we just want the max id
			if id > newID {
				newID = id
			}
		}

		newOptions := []validator.EnumVal{}
		// get custom options
		for _, curOption := range curOptions {
			exists := false
			for _, expectOption := range expectOptions {
				if expectOption.Name == curOption.Name {
					newOptions = append(newOptions, expectOption)
					exists = true
					break
				}
			}
			if !exists {
				if curOption.ID != "" {
					newOptions = append(newOptions, curOption)
					continue
				}
				newID++
				curOption.ID = strconv.Itoa(newID)
				newOptions = append(newOptions, curOption)
			}
		}

		blog.Infof("newOptions:= %#v", newOptions)
		if len(newOptions) <= 0 {
			continue
		}

		// update inst
		tablename := commondata.GetInstTableName(curRow.ObjectID)
		blog.Infof("updating option for table %s, property %s, option: %v", tablename, curRow.PropertyID, newOptions)
		defaultValue := ""
		validValues := []string{}
		for _, option := range newOptions {
			validValues = append(validValues, option.ID)
			if option.IsDefault {
				defaultValue = option.ID
			}
			updateinstdata := map[string]interface{}{
				curRow.PropertyID: option.ID,
			}
			updateinstcondition := map[string]interface{}{
				curRow.PropertyID: option.Name,
			}
			if tablename == common.BKTableNameBaseInst {
				updateinstcondition[common.BKObjIDField] = curRow.ObjectID
			}

			blog.Infof("update inst table %s, condition %#v, data %#v", tablename, updateinstcondition, updateinstdata)
			err = db.UpdateByCondition(tablename, updateinstdata, updateinstcondition)
			if err != nil {
				blog.Errorf("upgradeGlobalization update inst error: %v", err)
				return err
			}
		}

		// update default
		if defaultValue != "" {
			err = db.AddColumn(tablename, &storage.Column{Name: curRow.PropertyID, Ext: defaultValue})
			if err != nil {
				blog.Errorf("upgradeGlobalization update inst error: %v", err)
				return err
			}
			updateinstdata := map[string]interface{}{
				curRow.PropertyID: defaultValue,
			}
			updateinstcondition := map[string]interface{}{
				curRow.PropertyID: map[string]interface{}{
					"$nin": validValues,
				},
			}
			if tablename == common.BKTableNameBaseInst {
				updateinstcondition[common.BKObjIDField] = curRow.ObjectID
			}

			blog.Infof("update inst table %s, condition %#v, data %#v", tablename, updateinstcondition, updateinstdata)
			err = db.UpdateByCondition(tablename, updateinstdata, updateinstcondition)
			if err != nil {
				blog.Errorf("upgradeGlobalization update inst error: %v", err)
				return err
			}
		}

		// update property's option fields
		selector := map[string]interface{}{
			common.BKObjIDField:      curRow.ObjectID,
			common.BKPropertyIDField: curRow.PropertyID,
		}
		updatedata := map[string]interface{}{
			common.BKOptionField: newOptions,
		}
		blog.Infof("update attdes table %s, condition %#v, data %#v", common.BKTableNameObjAttDes, selector, updatedata)
		err = db.UpdateByCondition(common.BKTableNameObjAttDes, updatedata, selector)
		if err != nil {
			blog.Errorf("upgradeGlobalization update option field error: %v", err)
			return err
		}
	}
	return nil
}

func upgradeAppfield(instData storage.DI) error {
	condition := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDApp,
		common.BKPropertyIDField: map[string]interface{}{
			"$in": []string{
				"time_zone",
				"language",
			},
		},
	}
	data := map[string]interface{}{
		"isrequired": true,
	}
	instData.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)

	instData.AddColumn(common.BKTableNameBaseApp, &storage.Column{Name: common.BKTimeZoneField, Ext: "Asia/Shanghai"})

	return nil
}
