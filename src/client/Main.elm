module Main exposing (main)

-- Standard imports

import Http
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)


-- Local imports

import Mode.Scoreboard as Scoreboard exposing (Scoreboard(..))
import Mode.Message as Message exposing (Message)
import Mode.SubmitForm as SubmitForm exposing (SubmitFlagResponse, SubmitResponse)


-- MAIN


main : Platform.Program Basics.Never Model Msg
main =
    Html.program
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }



-- MODEL


type Notification
    = Error String
    | Success String
    | None


type ViewMode
    = ScoreboardView
    | SubmitForm
    | MessagesView


type Msg
    = SwitchMode ViewMode
    | ScoreboardRetrieved (Result Http.Error Scoreboard)
    | FlagSubmitted (Result Http.Error SubmitFlagResponse)
    | MessagesRetrieved (Result Http.Error (List Message))
    | GotInput SubmitForm.Input
    | DismissNotification


type alias Model =
    { scoreboard : Scoreboard
    , submitTokenInput : String
    , flagInput : String
    , messages : List Message
    , mode : ViewMode
    , notification : Notification
    }


init : ( Model, Cmd Msg )
init =
    ( Model (Scoreboard []) "" "" [] ScoreboardView None, Scoreboard.query ScoreboardRetrieved )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SwitchMode ScoreboardView ->
            ( { model | mode = ScoreboardView, notification = None }, Scoreboard.query ScoreboardRetrieved )

        SwitchMode MessagesView ->
            ( { model | mode = MessagesView, notification = None }, Message.query MessagesRetrieved )

        SwitchMode mode ->
            ( { model | mode = mode, notification = None }, Cmd.none )

        ScoreboardRetrieved (Ok scoreboard) ->
            ( { model | scoreboard = scoreboard }, Cmd.none )

        ScoreboardRetrieved _ ->
            ( model, Cmd.none )

        MessagesRetrieved (Ok messages) ->
            ( { model | messages = messages }, Cmd.none )

        MessagesRetrieved (Err err) ->
            ( model, Cmd.none )

        GotInput (SubmitForm.SubmissionToken input) ->
            ( { model | submitTokenInput = input }, Cmd.none )

        GotInput (SubmitForm.Flag input) ->
            ( { model | flagInput = input }, Cmd.none )

        GotInput SubmitForm.Submit ->
            ( model, SubmitForm.mutation model.submitTokenInput model.flagInput FlagSubmitted )

        DismissNotification ->
            ( { model | notification = None }, Cmd.none )

        FlagSubmitted (Err err) ->
            ( { model | notification = Error "Error submitting flag. Try again later." }, Cmd.none )

        FlagSubmitted (Ok { submitFlag }) ->
            let
                notification =
                    if submitFlag.correct then
                        Success "Congratulations! Your submission was accepted."
                    else
                        Error "Your submission was incorrect."
            in
                ( { notification = notification
                  , scoreboard = Scoreboard submitFlag.teams
                  , submitTokenInput = ""
                  , flagInput = ""
                  , mode = ScoreboardView
                  , messages = model.messages
                  }
                , Cmd.none
                )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ viewNav model
        , viewNotification model
        , viewMode model
        ]


viewNav : Model -> Html Msg
viewNav model =
    let
        navLinks =
            case model.mode of
                ScoreboardView ->
                    [ li [ class "active" ] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]

                MessagesView ->
                    [ li [] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [ class "active" ] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]

                SubmitForm ->
                    [ li [] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [ class "active" ] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]
    in
        nav []
            [ div [ class "mainContent nav-wrapper" ]
                [ a [ href "#", class "brand-logo" ] [ text "Scoreboard" ] --img [ src "/img/logo.png" ] [] ]
                , ul [ id "nav-mobile", class "right" ] navLinks
                ]
            ]


viewNotification : Model -> Html Msg
viewNotification model =
    let
        _ =
            Debug.log "Model is " model

        ( kind, displayMode, message ) =
            case model.notification of
                Error message ->
                    ( "Error", "", message )

                Success message ->
                    ( "Success", "", message )

                None ->
                    ( "", "none", "" )
    in
        div
            [ id "notification"
            , class "mainContent"
            , style [ ( "display", displayMode ) ]
            ]
            [ div [ class "card waves-effect waves-light" ]
                [ div [ class "card-content" ]
                    [ span [ class "card-title gray-text text-darken-4" ] [ text kind ]
                    , p [] [ text message ]
                    ]
                , div [ class "card-action" ]
                    [ a [ href "#", onClick DismissNotification ] [ text "Dismiss" ]
                    ]
                ]
            ]


viewMode : Model -> Html Msg
viewMode model =
    let
        viewContent =
            case model.mode of
                ScoreboardView ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Scoreboard" ]
                    , Scoreboard.view model.scoreboard
                    ]

                SubmitForm ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Submit a flag" ]
                    , SubmitForm.view GotInput
                    ]

                MessagesView ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Messages" ]
                    , Message.view model.messages
                    ]
    in
        div [ class "mainContent" ]
            [ div [ id "content", class "card waves-effect waves-light" ]
                [ div [ class "card-content" ] viewContent
                ]
            ]
