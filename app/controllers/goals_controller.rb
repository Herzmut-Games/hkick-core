class GoalsController < ApplicationController
  def create
    case params[:team]
    when 'red'
        ActionCable.server.broadcast 'goal_notifications_channel', message: 'red'
    when 'white'
      ActionCable.server.broadcast 'goal_notifications_channel', message: 'white'
    end
  end
end