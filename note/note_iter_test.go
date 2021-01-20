package note

import "testing"

/*
* @CreateTime: 2021/1/13 16:09
* @Author: JemmyHu <hujm20151021@gmail.com>
* @Description:
 */

func TestNoteIter_Next(t *testing.T) {
	n := GetNoteBook(testUserName)
	iter := NewNoteIter(GetDbHeadKey(), n.DB, testUserName, testPwd)
	for {
		note := iter.Next()
		if note == nil {
			break
		}
		t.Log(note.String())
	}
}
