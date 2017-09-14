class CreateLoginHistories < ActiveRecord::Migration[5.1]
  def change
    create_table :login_histories do |t|
      t.string    :username
      t.datetime  :login_at
    end
  end
end
