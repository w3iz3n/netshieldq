package dao

import (
	"github.com/jinzhu/gorm"
	"pnetshieldq/entity"
)

type FriendshipDao struct {
	db *gorm.DB
}

func NewFriendshipDao(db *gorm.DB) *FriendshipDao {
	return &FriendshipDao{db: db}
}

func (dao *FriendshipDao) AddFriend(friendship *entity.Friendship) error {
	return dao.db.Table("friendships").Create(friendship).Error
}

func (dao *FriendshipDao) CheckFriendship(userID1, userID2 int) (bool, error) {
	var count int64
	err := dao.db.Table("friendships").
		Where("((user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)) AND status = 'accepted'", userID1, userID2, userID2, userID1).
		Count(&count).Error
	return count > 0, err
}

func (dao *FriendshipDao) RemoveFriend(userID1, userID2 int) error {
	return dao.db.Table("friendships").
		Where("((user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?))", userID1, userID2, userID2, userID1).
		Delete(entity.Friendship{}).Error // Assuming you want to delete rows, not the whole record
}

func (dao *FriendshipDao) GetFriends(userID int) ([]int, error) {
	var friendIDs []int // 结果切片
	// 执行查询
	rows, err := dao.db.Table("friendships").
		Select("CASE WHEN user_id1 = ? THEN user_id2 ELSE user_id1 END", userID).
		Where("(user_id1 = ? OR user_id2 = ?) AND status = 'accepted'", userID, userID).
		Rows() // 获取 *sql.Rows 对象
	if err != nil {
		return nil, err // 查询错误处理
	}
	defer rows.Close() // 确保在函数返回前关闭 rows

	var friendID int
	for rows.Next() {
		if err := rows.Scan(&friendID); err != nil {
			return nil, err // 处理扫描错误
		}
		friendIDs = append(friendIDs, friendID) // 将每个提取的 ID 加入切片
	}

	return friendIDs, nil // 返回填充好的切片
}

func (dao *FriendshipDao) UpdateStatus(userID1 int, userID2 int, newStatus string) error {
	query := `UPDATE friendships SET status = ? WHERE (user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)`
	err := dao.db.Exec(query, newStatus, userID1, userID2, userID2, userID1)
	return err.Error
}
func (dao *FriendshipDao) GetFriendRequests(userID int) ([]entity.Friendship, error) {
	var friendRequests []entity.Friendship
	err := dao.db.Table("friendships").
		Where("user_id2 = ? ", userID).
		Order("created_at desc").
		Limit(10).
		Find(&friendRequests).Error
	return friendRequests, err
}
