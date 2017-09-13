class CreateUsers < ActiveRecord::Migration[5.1]
  def up
    create_table :users do |t|
      t.string  :username, unique: true
      t.string  :name
      t.string  :password
    end
  end

  def down
    drop_table :users
  end
end
