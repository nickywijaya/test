class AddAdditionalFieldsToUsers < ActiveRecord::Migration[5.1]
  def change
    add_column :users, :email, :string, unique: true
    add_column :users, :active, :boolean, default: true
  end
end
